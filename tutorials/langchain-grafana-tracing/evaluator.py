"""Behavioral evaluator for evolving the tracing helpers in app.py."""

from __future__ import annotations

import hashlib
import importlib.util
import os
import time
import uuid
from pathlib import Path

from openevolve.evaluation_result import EvaluationResult
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import SimpleSpanProcessor
from opentelemetry.sdk.trace.export.in_memory_span_exporter import InMemorySpanExporter
from opentelemetry.trace import StatusCode


class _PromptValue:
    """Minimal LangChain-compatible prompt value for isolated model tests."""

    def __init__(self, text: str) -> None:
        self._text = text

    def to_string(self) -> str:
        return self._text


def _load_program(program_path: str):
    module_name = f"tracing_candidate_{uuid.uuid4().hex}"
    spec = importlib.util.spec_from_file_location(module_name, program_path)
    if spec is None or spec.loader is None:
        raise ImportError(f"Unable to load candidate: {program_path}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def _test_candidate(program_path: str) -> tuple[dict[str, float], dict[str, str]]:
    module = _load_program(program_path)
    if not callable(getattr(module, "traced_step", None)):
        raise AttributeError("Candidate must define traced_step")
    if not callable(getattr(module, "fake_model", None)):
        raise AttributeError("Candidate must define fake_model")

    exporter = InMemorySpanExporter()
    provider = TracerProvider()
    provider.add_span_processor(SimpleSpanProcessor(exporter))
    tracer = provider.get_tracer("openevolve.langchain-evaluator")
    module.TRACER = tracer

    # Successful step: output, hierarchy, attributes, and duration.
    with tracer.start_as_current_span("test.root") as root:
        expected_parent = root.get_span_context().span_id
        started = time.perf_counter()
        result = module.traced_step("test.success", "test", lambda value: value + 1).invoke(1)
        overhead_ms = (time.perf_counter() - started) * 1000
    spans = list(exporter.get_finished_spans())
    success = next((span for span in spans if span.name == "test.success"), None)
    functional = float(result == 2)
    hierarchy = float(
        success is not None
        and success.parent is not None
        and success.parent.span_id == expected_parent
    )
    success_attributes = float(
        success is not None
        and success.attributes.get("gen_ai.operation.name") == "test"
        and success.attributes.get("langchain.step.name") == "test.success"
        and isinstance(success.attributes.get("step.duration_ms"), (int, float))
    )
    success_status = float(
        success is not None and success.status.status_code is StatusCode.OK
    )

    # Failure path: preserve exception behavior and emit an error span.
    exporter.clear()

    def fail(_: object) -> object:
        raise ValueError("private failure detail")

    reraised = False
    try:
        module.traced_step("test.failure", "test", fail).invoke(None)
    except ValueError as error:
        reraised = str(error) == "private failure detail"
    failure_spans = list(exporter.get_finished_spans())
    failure = next((span for span in failure_spans if span.name == "test.failure"), None)
    error_handling = float(
        reraised
        and failure is not None
        and failure.status.status_code is StatusCode.ERROR
        and any(event.name == "exception" for event in failure.events)
    )
    error_duration = float(
        failure is not None
        and isinstance(failure.attributes.get("step.duration_ms"), (int, float))
    )

    # Privacy path: raw prompt is opt-in; default stores only its digest.
    exporter.clear()
    old_trace_content = os.environ.pop("TRACE_CONTENT", None)
    question = "How does tracing preserve privacy?"
    try:
        with tracer.start_as_current_span("test.model"):
            answer = module.fake_model(
                _PromptValue(f"System: concise\nHuman: {question}")
            )
    finally:
        if old_trace_content is not None:
            os.environ["TRACE_CONTENT"] = old_trace_content
    model_span = next(
        (span for span in exporter.get_finished_spans() if span.name == "test.model"),
        None,
    )
    expected_hash = hashlib.sha256(question.encode()).hexdigest()[:12]
    privacy = float(
        model_span is not None
        and "gen_ai.prompt" not in model_span.attributes
        and model_span.attributes.get("gen_ai.prompt.sha256") == expected_hash
    )
    model_telemetry = float(
        model_span is not None
        and model_span.attributes.get("gen_ai.request.model")
        == "tutorial-fake-model"
        and model_span.attributes.get("gen_ai.usage.input_tokens", 0) > 0
        and model_span.attributes.get("gen_ai.usage.output_tokens", 0) > 0
        and question in answer
    )

    # End-to-end chain remains usable and emits the three tutorial steps.
    exporter.clear()
    chain_question = "What is a runnable?"
    with tracer.start_as_current_span("test.chain"):
        chain_answer = module.build_chain().invoke({"question": chain_question})
    chain_spans = list(exporter.get_finished_spans())
    names = {span.name for span in chain_spans}
    chain_behavior = float(
        chain_question in chain_answer
        and {"prompt.render", "model.generate", "output.parse"}.issubset(names)
    )

    metrics = {
        "functional_correctness": functional,
        "trace_hierarchy": hierarchy,
        "success_attributes": success_attributes,
        "success_status": success_status,
        "error_handling": error_handling,
        "error_duration": error_duration,
        "prompt_privacy": privacy,
        "model_telemetry": model_telemetry,
        "chain_behavior": chain_behavior,
    }
    metrics["combined_score"] = float(
        0.10 * functional
        + 0.10 * hierarchy
        + 0.10 * success_attributes
        + 0.10 * success_status
        + 0.15 * error_handling
        + 0.15 * error_duration
        + 0.10 * privacy
        + 0.10 * model_telemetry
        + 0.10 * chain_behavior
    )
    artifacts = {
        "summary": f"Passed {sum(metrics[key] for key in metrics if key != 'combined_score'):.0f}/9 checks",
        "instrumentation_overhead_ms": f"{overhead_ms:.4f}",
    }
    provider.shutdown()
    return metrics, artifacts


def evaluate(program_path: str) -> EvaluationResult:
    """Score functional, tracing, error, and privacy invariants."""
    try:
        metrics, artifacts = _test_candidate(program_path)
        return EvaluationResult(metrics=metrics, artifacts=artifacts)
    except Exception as error:
        return EvaluationResult(
            metrics={"combined_score": 0.0},
            artifacts={"error": f"{type(error).__name__}: {error}"},
        )


if __name__ == "__main__":
    candidate = str(Path(__file__).with_name("app.py"))
    result = evaluate(candidate)
    print(result.metrics)
    print(result.artifacts)

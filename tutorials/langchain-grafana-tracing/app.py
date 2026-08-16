"""Small LangChain pipeline with vendor-neutral OpenTelemetry tracing."""

from __future__ import annotations

import argparse
import hashlib
import os
import time
from collections.abc import Callable
from typing import Any

from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.runnables import RunnableLambda
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace import Status, StatusCode


def configure_tracing() -> TracerProvider:
    """Configure batching once at the application boundary."""
    resource = Resource.create(
        {
            "service.name": os.getenv("OTEL_SERVICE_NAME", "langchain-tutorial"),
            "service.version": os.getenv("SERVICE_VERSION", "0.1.0"),
            "deployment.environment.name": os.getenv(
                "DEPLOYMENT_ENVIRONMENT", "local"
            ),
        }
    )
    provider = TracerProvider(resource=resource)
    provider.add_span_processor(
        BatchSpanProcessor(
            OTLPSpanExporter(
                endpoint=os.getenv(
                    "OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"
                ),
                insecure=True,
            )
        )
    )
    trace.set_tracer_provider(provider)
    return provider


TRACER = trace.get_tracer("tutorial.langchain")


# EVOLVE-BLOCK-START
def traced_step(name: str, operation: str, function: Callable[[Any], Any]):
    """Wrap one LangChain runnable in a child span."""

    def invoke(value: Any) -> Any:
        started = time.perf_counter()
        with TRACER.start_as_current_span(name) as span:
            span.set_attribute("gen_ai.operation.name", operation)
            span.set_attribute("langchain.step.name", name)
            try:
                result = function(value)
                span.set_attribute(
                    "step.duration_ms", (time.perf_counter() - started) * 1000
                )
                return result
            except Exception as error:
                span.record_exception(error)
                span.set_status(Status(StatusCode.ERROR, str(error)))
                raise

    return RunnableLambda(invoke)


def fake_model(prompt: Any) -> str:
    """A free deterministic model substitute, so the lab needs no API key."""
    time.sleep(0.08)
    text = prompt.to_string()
    question = text.rsplit("Human:", maxsplit=1)[-1].strip()
    if os.getenv("TRACE_CONTENT", "false").lower() == "true":
        trace.get_current_span().set_attribute("gen_ai.prompt", question[:1000])
    else:
        digest = hashlib.sha256(question.encode()).hexdigest()[:12]
        trace.get_current_span().set_attribute("gen_ai.prompt.sha256", digest)
    trace.get_current_span().set_attribute("gen_ai.request.model", "tutorial-fake-model")
    trace.get_current_span().set_attribute("gen_ai.usage.input_tokens", len(text.split()))
    answer = f"A concise tutorial answer for: {question}"
    trace.get_current_span().set_attribute("gen_ai.usage.output_tokens", len(answer.split()))
    return answer
# EVOLVE-BLOCK-END


def build_chain():
    prompt = ChatPromptTemplate.from_messages(
        [
            ("system", "You teach LangChain in short, concrete steps."),
            ("human", "{question}"),
        ]
    )
    return (
        traced_step("prompt.render", "prompt", prompt.invoke)
        | traced_step("model.generate", "chat", fake_model)
        | traced_step("output.parse", "parse", StrOutputParser().invoke)
    )


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "question", nargs="?", default="What problem does LangChain solve?"
    )
    args = parser.parse_args()
    provider = configure_tracing()
    try:
        with TRACER.start_as_current_span("langchain.request") as span:
            span.set_attribute("app.request.type", "tutorial_question")
            print(build_chain().invoke({"question": args.question}))
            print(f"trace_id={span.get_span_context().trace_id:032x}")
    finally:
        provider.shutdown()


if __name__ == "__main__":
    main()

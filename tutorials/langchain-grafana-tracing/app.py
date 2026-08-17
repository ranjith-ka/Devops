"""CLI orchestration and main entry point.

This module handles command-line argument parsing and workflow orchestration
for the LangChain tracing and trace comparison application.
"""

from __future__ import annotations

import argparse
import json

from graph import AgentGraph
from tempo import fetch_trace_from_tempo
from trace_analyzer import compare_traces
from tracing import configure_tracing, get_tracer
from opentelemetry import trace

TRACER = None


def run_question_answering(question: str, thread_id: str = "cli") -> None:
    """Run the LangChain pipeline for question answering."""
    global TRACER
    if TRACER is None:
        configure_tracing()
        TRACER = get_tracer()

    with TRACER.start_as_current_span("langchain.request") as span:
        span.set_attribute("app.request.type", "tutorial_question")
        graph = AgentGraph()
        try:
            result = graph.invoke(question, thread_id)
        finally:
            graph.close()
        print(result["answer"])
        print(f"trace_id={span.get_span_context().trace_id:032x}")

    provider = trace.get_tracer_provider()
    if hasattr(provider, "force_flush"):
        provider.force_flush()


def run_trace_comparison(trace_a_id: str, trace_b_id: str) -> None:
    """Run the trace comparison workflow."""
    trace_a = fetch_trace_from_tempo(trace_a_id)
    trace_b = fetch_trace_from_tempo(trace_b_id)
    result = compare_traces(trace_a, trace_b)

    print("\n=== Trace Comparison Report ===\n")
    print(json.dumps(result.to_dict(), indent=2))
    print(f"\nRoot Cause: {result.root_cause}")
    print(f"Recommendation: {result.recommendation}")


def main() -> None:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="LangChain tracing and trace comparison demo"
    )
    parser.add_argument(
        "question",
        nargs="?",
        default="What problem does LangChain solve?",
        help="Question to ask the model",
    )
    parser.add_argument(
        "--compare-trace-a",
        help="First trace ID for comparison",
    )
    parser.add_argument(
        "--compare-trace-b",
        help="Second trace ID for comparison",
    )
    parser.add_argument("--thread-id", default="cli", help="Conversation memory thread ID")

    args = parser.parse_args()
    provider = configure_tracing()
    global TRACER
    TRACER = get_tracer()

    try:
        if args.compare_trace_a and args.compare_trace_b:
            run_trace_comparison(args.compare_trace_a, args.compare_trace_b)
        else:
            run_question_answering(args.question, args.thread_id)
    finally:
        provider.force_flush()
        provider.shutdown()


if __name__ == "__main__":
    main()

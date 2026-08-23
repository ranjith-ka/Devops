"""Read-only LangChain tools for investigating Tempo traces and Loki logs."""

from __future__ import annotations

import json
import time
from collections.abc import Callable
from typing import Any

from langchain_core.tools import tool
from opentelemetry.trace import Status, StatusCode

from loki import emit_log, fetch_logs_for_trace
from tempo import fetch_trace_from_tempo
from trace_analyzer import compare_traces
from tracing import get_tracer

MAX_TRACE_SPANS = 40
MAX_LOG_RECORDS = 50


def _trace_id(value: str) -> str:
    normalized = value.strip().lower()
    if len(normalized) != 32 or any(character not in "0123456789abcdef" for character in normalized):
        raise ValueError("Trace ID must be a 32-character hexadecimal value")
    return normalized


def _safe_json_tool(name: str, operation: Callable[[], dict[str, Any]]) -> str:
    """Run a tool with bounded, observable, model-readable error handling."""
    started = time.perf_counter()
    tracer = get_tracer("langchain.incident_agent")
    with tracer.start_as_current_span(f"agent.tool.{name}") as span:
        span.set_attribute("gen_ai.operation.name", "execute_tool")
        span.set_attribute("gen_ai.tool.name", name)
        emit_log("agent.tool.started", tool=name)
        try:
            result = operation()
            payload = json.dumps(result, default=str, separators=(",", ":"))
            span.set_attribute("tool.result.bytes", len(payload.encode("utf-8")))
            span.set_attribute("tool.status", "ok")
            emit_log("agent.tool.completed", tool=name, result_bytes=len(payload.encode("utf-8")))
            return payload
        except Exception as error:
            span.record_exception(error)
            span.set_status(Status(StatusCode.ERROR, type(error).__name__))
            span.set_attribute("tool.status", "error")
            emit_log("agent.tool.failed", level="error", tool=name, error_type=type(error).__name__)
            return json.dumps({
                "ok": False,
                "tool": name,
                "error_type": type(error).__name__,
                "message": str(error),
            }, separators=(",", ":"))


@tool
def get_trace(trace_id: str) -> str:
    """Fetch one Tempo trace by its 32-character hex ID and return compact span evidence."""
    def operation() -> dict[str, Any]:
        trace = fetch_trace_from_tempo(_trace_id(trace_id))
        ordered = sorted(trace.spans, key=lambda span: span.duration_ms, reverse=True)
        return {
            "ok": True,
            "trace_id": trace.trace_id,
            "root": {
                "name": trace.root_name,
                "duration_ms": round(trace.root_duration_ms, 3),
                "status": trace.root_status,
            },
            "span_count": len(trace.spans),
            "spans_truncated": len(trace.spans) > MAX_TRACE_SPANS,
            "slowest_spans": [
                {
                    "name": span.name,
                    "duration_ms": round(span.duration_ms, 3),
                    "status": span.status,
                    "start_offset_ms": round(span.start_offset_ms, 3),
                }
                for span in ordered[:MAX_TRACE_SPANS]
            ],
        }

    return _safe_json_tool("get_trace", operation)


@tool
def get_logs(trace_id: str, limit: int = 30) -> str:
    """Fetch Loki application logs correlated with a trace ID when trace evidence needs deeper investigation."""
    def operation() -> dict[str, Any]:
        normalized_id = _trace_id(trace_id)
        bounded_limit = max(1, min(int(limit), MAX_LOG_RECORDS))
        records = fetch_logs_for_trace(normalized_id, limit=bounded_limit)
        compact = []
        for record in records[:bounded_limit]:
            compact.append({
                key: value
                for key, value in record.items()
                if key in {"timestamp", "level", "event", "trace_id", "node", "error_type", "tool"}
            })
        return {
            "ok": True,
            "trace_id": normalized_id,
            "record_count": len(records),
            "records_truncated": len(records) > bounded_limit,
            "records": compact,
        }

    return _safe_json_tool("get_logs", operation)


@tool
def compare_trace_ids(baseline_trace_id: str, candidate_trace_id: str) -> str:
    """Compare a baseline Tempo trace with a candidate trace to find structural or latency regressions."""
    def operation() -> dict[str, Any]:
        baseline = fetch_trace_from_tempo(_trace_id(baseline_trace_id))
        candidate = fetch_trace_from_tempo(_trace_id(candidate_trace_id))
        comparison = compare_traces(baseline, candidate)
        ranked = sorted(
            comparison.span_details.values(),
            key=lambda item: abs(item.delta_ms),
            reverse=True,
        )[:20]
        return {
            "ok": True,
            "baseline_trace_id": comparison.trace_a_id,
            "candidate_trace_id": comparison.trace_b_id,
            "root_diff_ms": round(comparison.root_diff_ms, 3),
            "slowest_span": comparison.slowest_span,
            "root_cause_hypothesis": comparison.root_cause,
            "recommendation": comparison.recommendation,
            "largest_span_deltas": [
                {
                    "name": item.name,
                    "baseline_ms": round(item.trace_a_ms, 3),
                    "candidate_ms": round(item.trace_b_ms, 3),
                    "delta_ms": round(item.delta_ms, 3),
                    "baseline_status": item.status_a,
                    "candidate_status": item.status_b,
                }
                for item in ranked
            ],
        }

    return _safe_json_tool("compare_trace_ids", operation)


INCIDENT_TOOLS = [get_trace, get_logs, compare_trace_ids]

"""Fetch and normalize traces from Grafana Tempo."""

from __future__ import annotations

import json
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.parse import quote
from urllib.request import Request, urlopen

from config import Config
from trace_data import TraceData, TraceSpan


def _any_value(value: dict[str, Any]) -> Any:
    """Convert an OTLP JSON AnyValue into a regular Python value."""
    for key in ("stringValue", "boolValue", "intValue", "doubleValue", "bytesValue"):
        if key in value:
            return value[key]
    if "arrayValue" in value:
        return [_any_value(item) for item in value["arrayValue"].get("values", [])]
    if "kvlistValue" in value:
        return _attributes(value["kvlistValue"].get("values", []))
    return None


def _attributes(attributes: list[dict[str, Any]]) -> dict[str, Any]:
    return {
        attribute["key"]: _any_value(attribute.get("value", {}))
        for attribute in attributes
        if "key" in attribute
    }


def _span_status(span: dict[str, Any]) -> str:
    code = str(span.get("status", {}).get("code", "")).upper()
    return "error" if code in {"2", "STATUS_CODE_ERROR"} else "ok"


def _duration_ns(span: dict[str, Any]) -> int:
    try:
        start = int(span.get("startTimeUnixNano", 0))
        end = int(span.get("endTimeUnixNano", 0))
    except (TypeError, ValueError):
        return 0
    return max(0, end - start)


def parse_tempo_trace(trace_id: str, payload: dict[str, Any]) -> TraceData:
    """Convert Tempo's OTLP JSON response into the comparison data model."""
    spans: list[dict[str, Any]] = []
    for batch in payload.get("batches", []):
        resource_attributes = _attributes(batch.get("resource", {}).get("attributes", []))
        for scope_spans in batch.get("scopeSpans", []):
            for span in scope_spans.get("spans", []):
                normalized = dict(span)
                normalized["_resource_attributes"] = resource_attributes
                spans.append(normalized)

    if not spans:
        raise ValueError(f"Tempo returned no spans for trace {trace_id}")

    root_candidates = [span for span in spans if not span.get("parentSpanId")]
    root = max(root_candidates or spans, key=_duration_ns)
    root_start_ns = int(root.get("startTimeUnixNano", 0))
    child_spans = [
        TraceSpan(
            name=span.get("name", "unnamed"),
            duration_ms=_duration_ns(span) / 1_000_000,
            status=_span_status(span),
            start_offset_ms=max(
                0, (int(span.get("startTimeUnixNano", root_start_ns)) - root_start_ns) / 1_000_000
            ),
            metadata={
                **span.get("_resource_attributes", {}),
                **_attributes(span.get("attributes", [])),
            },
        )
        for span in spans
        if span is not root
    ]

    return TraceData(
        trace_id=trace_id,
        root_duration_ms=_duration_ns(root) / 1_000_000,
        root_status=_span_status(root),
        spans=child_spans,
        root_name=root.get("name", "root"),
    )


def fetch_trace_from_tempo(trace_id: str) -> TraceData:
    """Fetch a real trace from Tempo by its 32-character hexadecimal ID."""
    normalized_id = trace_id.strip().lower()
    if len(normalized_id) != 32 or any(char not in "0123456789abcdef" for char in normalized_id):
        raise ValueError("Trace ID must be a 32-character hexadecimal value")

    url = f"{Config.TEMPO_ENDPOINT}/api/traces/{quote(normalized_id, safe='')}"
    request = Request(url, headers={"Accept": "application/json"})
    try:
        with urlopen(request, timeout=10) as response:
            payload = json.load(response)
    except HTTPError as error:
        if error.code == 404:
            raise ValueError(f"Trace {normalized_id} was not found in Tempo") from error
        raise RuntimeError(f"Tempo returned HTTP {error.code} for trace {normalized_id}") from error
    except URLError as error:
        raise RuntimeError(f"Could not connect to Tempo at {Config.TEMPO_ENDPOINT}: {error.reason}") from error

    return parse_tempo_trace(normalized_id, payload)

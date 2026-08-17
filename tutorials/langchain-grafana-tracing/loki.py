"""Push and query trace-correlated structured application logs in Loki."""

from __future__ import annotations

import json
import time
from typing import Any
from urllib.error import HTTPError, URLError
from urllib.parse import urlencode
from urllib.request import Request, urlopen

from opentelemetry import trace

from config import Config


def current_trace_id() -> str:
    context = trace.get_current_span().get_span_context()
    return f"{context.trace_id:032x}" if context.is_valid else ""


def emit_log(event: str, *, level: str = "info", **fields: Any) -> None:
    """Best-effort structured log delivery; observability must not break the app."""
    record = {
        "timestamp": time.time_ns(),
        "level": level,
        "event": event,
        "trace_id": current_trace_id(),
        **fields,
    }
    payload = {
        "streams": [{
            "stream": {
                "service_name": Config.SERVICE_NAME,
                "environment": Config.ENVIRONMENT,
            },
            "values": [[str(record["timestamp"]), json.dumps(record, default=str)]],
        }]
    }
    request = Request(
        f"{Config.LOKI_ENDPOINT}/loki/api/v1/push",
        data=json.dumps(payload).encode("utf-8"),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urlopen(request, timeout=2):
            pass
    except (HTTPError, URLError, TimeoutError, OSError):
        return


def fetch_logs_for_trace(trace_id: str, limit: int = 200) -> list[dict[str, Any]]:
    """Return structured Loki records containing the requested trace ID."""
    query = f'{{service_name="{Config.SERVICE_NAME}"}} |= "{trace_id}"'
    params = urlencode({"query": query, "limit": limit, "direction": "forward"})
    request = Request(
        f"{Config.LOKI_ENDPOINT}/loki/api/v1/query_range?{params}",
        headers={"Accept": "application/json"},
    )
    try:
        with urlopen(request, timeout=5) as response:
            payload = json.load(response)
    except (HTTPError, URLError, TimeoutError, OSError, json.JSONDecodeError):
        return []

    records: list[dict[str, Any]] = []
    for stream in payload.get("data", {}).get("result", []):
        for timestamp, line in stream.get("values", []):
            try:
                record = json.loads(line)
            except json.JSONDecodeError:
                record = {"event": line}
            record["timestamp"] = int(timestamp)
            records.append(record)
    return sorted(records, key=lambda record: record.get("timestamp", 0))

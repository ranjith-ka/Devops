"""Data models for trace representation."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


@dataclass
class TraceSpan:
    """Represents a single span in a trace."""

    name: str
    duration_ms: float
    status: str
    start_offset_ms: float = 0
    metadata: dict[str, Any] = field(default_factory=dict)


@dataclass
class TraceData:
    """Represents a complete trace with root and child spans."""

    trace_id: str
    root_duration_ms: float
    root_status: str
    spans: list[TraceSpan]
    root_name: str = "root"

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> TraceData:
        """Convert a dict to TraceData."""
        spans = [
            TraceSpan(
                name=s["name"],
                duration_ms=s.get("duration_ms", 0),
                status=s.get("status", "ok"),
                start_offset_ms=s.get("start_offset_ms", 0),
                metadata=s.get("metadata", {}),
            )
            for s in data.get("children", [])
        ]
        root = data.get("root", {})
        return cls(
            trace_id=data["trace_id"],
            root_duration_ms=root.get("duration_ms", 0),
            root_status=root.get("status", "ok"),
            spans=spans,
            root_name=root.get("name", "root"),
        )


@dataclass
class SpanComparison:
    """Comparison of a span across two traces."""

    name: str
    trace_a_ms: float
    trace_b_ms: float
    delta_ms: float
    status_a: str
    status_b: str


@dataclass
class TraceComparisonResult:
    """Result of comparing two traces."""

    trace_a_id: str
    trace_b_id: str
    root_diff_ms: float
    slowest_span: str
    root_cause: str
    recommendation: str
    span_details: dict[str, SpanComparison]
    trace_a_waterfall: dict[str, Any]
    trace_b_waterfall: dict[str, Any]

    def to_dict(self) -> dict[str, Any]:
        """Convert to a serializable dict."""
        return {
            "trace_a": self.trace_a_id,
            "trace_b": self.trace_b_id,
            "root_diff_ms": self.root_diff_ms,
            "slowest_span": self.slowest_span,
            "root_cause": self.root_cause,
            "recommendation": self.recommendation,
            "trace_a_waterfall": self.trace_a_waterfall,
            "trace_b_waterfall": self.trace_b_waterfall,
            "span_details": {
                name: {
                    "trace_a_ms": comp.trace_a_ms,
                    "trace_b_ms": comp.trace_b_ms,
                    "delta_ms": comp.delta_ms,
                    "status_a": comp.status_a,
                    "status_b": comp.status_b,
                }
                for name, comp in self.span_details.items()
            },
        }

"""Trace comparison and root cause analysis."""

from trace_data import TraceData, TraceComparisonResult, SpanComparison, TraceSpan


def _waterfall(trace: TraceData) -> dict:
    """Serialize one trace as ordered waterfall lanes."""
    spans = [{
        "name": trace.root_name,
        "start_offset_ms": 0,
        "duration_ms": trace.root_duration_ms,
        "status": trace.root_status,
        "is_root": True,
    }]
    spans.extend({
        "name": span.name,
        "start_offset_ms": span.start_offset_ms,
        "duration_ms": span.duration_ms,
        "status": span.status,
        "is_root": False,
    } for span in sorted(trace.spans, key=lambda item: item.start_offset_ms))
    return {"trace_id": trace.trace_id, "duration_ms": trace.root_duration_ms, "spans": spans}


def diagnose_root_cause(
    trace_a: TraceData, trace_b: TraceData, details: dict[str, SpanComparison]
) -> tuple[str, str]:
    """Diagnose the root cause of trace differences."""
    if trace_a.root_status != trace_b.root_status:
        return (
            "status mismatch or validation error path",
            "Check guardrail logic, validation step, and any conditional retry path.",
        )

    missing = [name for name, comp in details.items() if "missing" in {comp.status_a, comp.status_b}]
    if missing:
        return (
            f"Trace structure differs: {', '.join(missing)}",
            "Check conditional branches, retries, and instrumentation coverage.",
        )

    meaningful_regressions = []
    for span_name, comp in details.items():
        threshold_ms = max(25.0, comp.trace_a_ms * 0.10)
        if comp.delta_ms >= threshold_ms:
            meaningful_regressions.append((span_name, comp))

    if meaningful_regressions:
        span_name, comp = max(meaningful_regressions, key=lambda item: item[1].delta_ms)
        percent = (comp.delta_ms / comp.trace_a_ms * 100) if comp.trace_a_ms else 0
        base_name = span_name.split(" [", 1)[0]
        label = base_name.replace("graph.node.", "").replace(".", " ").replace("_", " ").title()
        if base_name in {"model.generate", "graph.node.generate_answer"}:
            recommendation = "Check provider load, token size, model fallback, and timeout config."
        elif "tool" in base_name:
            recommendation = "Check external dependency latency, retry policy, and tool timeouts."
        elif "retriev" in base_name or "documentation" in base_name:
            recommendation = "Inspect document/vector lookup latency and index/query efficiency."
        else:
            recommendation = f"Inspect the {base_name} node inputs, dependencies, retries, and logs."
        return (
            f"{label} is {comp.delta_ms:.1f} ms slower ({percent:.1f}%)",
            recommendation,
        )

    return (
        "No major root cause found; traces are similar",
        "No corrective action required based on span comparison.",
    )


def compare_traces(trace_a: TraceData, trace_b: TraceData) -> TraceComparisonResult:
    """Compare two traces and identify root causes."""
    # Build span comparison details
    def keyed_spans(trace: TraceData) -> dict[str, TraceSpan]:
        counts: dict[str, int] = {}
        keyed = {}
        for span in sorted(trace.spans, key=lambda item: item.start_offset_ms):
            counts[span.name] = counts.get(span.name, 0) + 1
            occurrence = counts[span.name]
            key = span.name if occurrence == 1 else f"{span.name} [{occurrence}]"
            keyed[key] = span
        return keyed

    spans_a = keyed_spans(trace_a)
    spans_b = keyed_spans(trace_b)
    span_names = set(spans_a) | set(spans_b)
    details: dict[str, SpanComparison] = {}

    for span_name in sorted(span_names):
        span_a = spans_a.get(span_name)
        span_b = spans_b.get(span_name)

        a_ms = span_a.duration_ms if span_a else 0
        b_ms = span_b.duration_ms if span_b else 0

        details[span_name] = SpanComparison(
            name=span_name,
            trace_a_ms=a_ms,
            trace_b_ms=b_ms,
            delta_ms=b_ms - a_ms,
            status_a=span_a.status if span_a else "missing",
            status_b=span_b.status if span_b else "missing",
        )

    # Find slowest span
    slowest = max(
        details.items(),
        key=lambda item: abs(item[1].delta_ms) + max(item[1].trace_a_ms, item[1].trace_b_ms),
        default=(None, None),
    )[0] or "root"

    # Diagnose root cause
    root_cause, recommendation = diagnose_root_cause(trace_a, trace_b, details)

    return TraceComparisonResult(
        trace_a_id=trace_a.trace_id,
        trace_b_id=trace_b.trace_id,
        root_diff_ms=trace_b.root_duration_ms - trace_a.root_duration_ms,
        slowest_span=slowest,
        root_cause=root_cause,
        recommendation=recommendation,
        span_details=details,
        trace_a_waterfall=_waterfall(trace_a),
        trace_b_waterfall=_waterfall(trace_b),
    )

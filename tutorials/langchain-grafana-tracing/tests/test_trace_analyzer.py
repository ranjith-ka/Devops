import unittest

from trace_analyzer import compare_traces
from trace_data import TraceData, TraceSpan


def make_trace(trace_id, root_ms, model_ms):
    return TraceData(trace_id, root_ms, "ok", [TraceSpan("model.generate", model_ms, "ok", 50)])


class TraceAnalyzerTest(unittest.TestCase):
    def test_ignores_small_model_jitter(self):
        result = compare_traces(make_trace("a", 1100, 1000), make_trace("b", 1120, 1020))
        self.assertIn("similar", result.root_cause)

    def test_reports_measured_model_regression(self):
        result = compare_traces(make_trace("a", 1100, 1000), make_trace("b", 1400, 1300))
        self.assertEqual(result.root_cause, "Model Generate is 300.0 ms slower (30.0%)")

    def test_compares_repeated_span_occurrences(self):
        trace_a = TraceData("a", 100, "ok", [
            TraceSpan("tool", 10, "ok", 5), TraceSpan("tool", 20, "ok", 20)
        ])
        trace_b = TraceData("b", 100, "ok", [
            TraceSpan("tool", 10, "ok", 5), TraceSpan("tool", 60, "ok", 20)
        ])
        result = compare_traces(trace_a, trace_b)
        self.assertIn("tool [2]", result.span_details)
        self.assertEqual(result.span_details["tool [2]"].delta_ms, 40)


if __name__ == "__main__":
    unittest.main()

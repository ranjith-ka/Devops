import unittest

from tempo import parse_tempo_trace


class TempoTraceParsingTest(unittest.TestCase):
    def test_parses_real_otlp_tempo_response(self):
        payload = {
            "batches": [{
                "resource": {"attributes": [{
                    "key": "service.name", "value": {"stringValue": "demo-service"}
                }]},
                "scopeSpans": [{"spans": [
                    {
                        "name": "langchain.request", "spanId": "root",
                        "startTimeUnixNano": "1000000000", "endTimeUnixNano": "4000000000",
                        "status": {},
                    },
                    {
                        "name": "model.generate", "spanId": "child", "parentSpanId": "root",
                        "startTimeUnixNano": "1500000000", "endTimeUnixNano": "3500000000",
                        "status": {"code": "STATUS_CODE_ERROR"},
                        "attributes": [{
                            "key": "gen_ai.request.model", "value": {"stringValue": "llama3.2"}
                        }],
                    },
                ]}],
            }]}

        trace = parse_tempo_trace("a" * 32, payload)

        self.assertEqual(trace.root_duration_ms, 3000)
        self.assertEqual(trace.root_status, "ok")
        self.assertEqual(len(trace.spans), 1)
        self.assertEqual(trace.spans[0].duration_ms, 2000)
        self.assertEqual(trace.spans[0].start_offset_ms, 500)
        self.assertEqual(trace.spans[0].status, "error")
        self.assertEqual(trace.spans[0].metadata["service.name"], "demo-service")
        self.assertEqual(trace.spans[0].metadata["gen_ai.request.model"], "llama3.2")

    def test_rejects_empty_trace(self):
        with self.assertRaisesRegex(ValueError, "no spans"):
            parse_tempo_trace("a" * 32, {"batches": []})


if __name__ == "__main__":
    unittest.main()

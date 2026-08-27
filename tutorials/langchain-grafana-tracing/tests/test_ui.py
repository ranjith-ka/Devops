"""HTTP contract tests for the TracePilot Flask UI."""

import json
import unittest
from unittest.mock import patch

from langchain_core.messages import AIMessage, ToolMessage

import ui


class TracePilotUITests(unittest.TestCase):
    def setUp(self):
        ui.app.config.update(TESTING=True)
        self.client = ui.app.test_client()

    def test_home_renders_conversational_workspace(self):
        response = self.client.get("/")
        self.assertEqual(response.status_code, 200)
        self.assertIn(b"Turn telemetry into", response.data)
        self.assertIn(b"Compare two traces", response.data)
        self.assertEqual(response.headers["Cache-Control"], "no-store, max-age=0")

    def test_question_rejects_missing_text(self):
        response = self.client.post("/api/question", json={"question": "  "})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.get_json()["code"], "invalid_question")

    def test_compare_rejects_malformed_trace_ids(self):
        response = self.client.post("/api/compare", json={"trace_a": "abc", "trace_b": "def"})
        self.assertEqual(response.status_code, 400)
        self.assertEqual(response.get_json()["code"], "invalid_trace_ids")

    def test_status_reports_each_dependency(self):
        with patch.object(ui, "_service_status", side_effect=lambda name, _url: {"name": name, "available": True}):
            response = self.client.get("/api/status")
        self.assertEqual(response.status_code, 200)
        self.assertEqual([item["name"] for item in response.get_json()["services"]], ["Tempo", "Loki", "Ollama"])

    def test_incident_messages_become_structured_evidence(self):
        messages = [
            AIMessage(content="", tool_calls=[{"name": "get_trace", "args": {}, "id": "call-1"}]),
            ToolMessage(content=json.dumps({"ok": True, "span_count": 4}), tool_call_id="call-1", name="get_trace"),
            AIMessage(content="Evidence: four spans. Confidence: high. Next action: inspect the slowest span."),
        ]
        payload = ui._incident_payload({"messages": messages})
        self.assertEqual(payload["selected_tools"], ["get_trace"])
        self.assertEqual(payload["evidence"][0]["data"]["span_count"], 4)
        self.assertIn("Confidence", payload["answer"])


if __name__ == "__main__":
    unittest.main()

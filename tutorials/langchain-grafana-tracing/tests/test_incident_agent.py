import tempfile
import unittest
import json
from pathlib import Path

from langchain_core.messages import AIMessage, ToolMessage
from langchain_core.tools import tool

from incident_agent import IncidentAgent


@tool
def fake_trace(trace_id: str) -> str:
    """Return deterministic trace evidence."""
    return json.dumps({"trace_id": trace_id, "status": "error", "slow_service": "payment"})


@tool
def fake_logs(trace_id: str) -> str:
    """Return deterministic log evidence."""
    return json.dumps({"trace_id": trace_id, "event": "connection_pool_timeout"})


class ScriptedModel:
    def __init__(self, steps):
        self.steps = list(steps)
        self.calls = 0

    def bind_tools(self, tools):
        self.tools = tools
        return self

    def invoke(self, messages):
        step = self.steps[min(self.calls, len(self.steps) - 1)]
        self.calls += 1
        return step


class IncidentAgentTest(unittest.TestCase):
    def make_agent(self, model, directory):
        return IncidentAgent(
            model=model,
            tools=[fake_trace, fake_logs],
            checkpoint_path=str(Path(directory, "incident.sqlite")),
        )

    def test_agent_stops_after_trace_when_evidence_is_sufficient(self):
        trace_id = "a" * 32
        model = ScriptedModel([
            AIMessage(content="", tool_calls=[{
                "name": "fake_trace", "args": {"trace_id": trace_id}, "id": "trace-1"
            }]),
            AIMessage(content="Evidence: healthy enough. Confidence: high. Next action: none."),
        ])
        with tempfile.TemporaryDirectory() as directory:
            agent = self.make_agent(model, directory)
            try:
                result = agent.invoke(f"Investigate {trace_id}", "healthy")
            finally:
                agent.close()

        tool_names = [message.name for message in result["messages"] if isinstance(message, ToolMessage)]
        self.assertEqual(tool_names, ["fake_trace"])
        self.assertEqual(model.calls, 2)

    def test_agent_investigates_deeper_when_trace_requires_logs(self):
        trace_id = "b" * 32
        model = ScriptedModel([
            AIMessage(content="", tool_calls=[{
                "name": "fake_trace", "args": {"trace_id": trace_id}, "id": "trace-1"
            }]),
            AIMessage(content="", tool_calls=[{
                "name": "fake_logs", "args": {"trace_id": trace_id}, "id": "logs-1"
            }]),
            AIMessage(content="Evidence: pool timeout. Root-cause hypothesis: exhausted pool. Confidence: medium. Next action: inspect pool metrics."),
        ])
        with tempfile.TemporaryDirectory() as directory:
            agent = self.make_agent(model, directory)
            try:
                result = agent.invoke(f"Deeply investigate {trace_id}", "error")
            finally:
                agent.close()

        tool_names = [message.name for message in result["messages"] if isinstance(message, ToolMessage)]
        self.assertEqual(tool_names, ["fake_trace", "fake_logs"])
        self.assertEqual(model.calls, 3)

    def test_turn_limit_stops_repeated_tool_calls(self):
        trace_id = "c" * 32
        repeated = AIMessage(content="", tool_calls=[{
            "name": "fake_trace", "args": {"trace_id": trace_id}, "id": "repeat"
        }])
        model = ScriptedModel([repeated])
        with tempfile.TemporaryDirectory() as directory:
            agent = self.make_agent(model, directory)
            try:
                result = agent.invoke(f"Investigate {trace_id}", "bounded", max_model_calls=2)
            finally:
                agent.close()

        self.assertIn("model-turn limit", result["messages"][-1].content)
        self.assertEqual(model.calls, 2)


if __name__ == "__main__":
    unittest.main()

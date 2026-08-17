import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from langchain_core.messages import AIMessage

from config import Config
from documentation import DocumentationStore
from graph import AgentGraph


class FakeModel:
    def invoke(self, messages):
        return AIMessage(content=messages[-1].content)


class AgentGraphTest(unittest.TestCase):
    def test_documentation_retrieval_and_thread_memory(self):
        with tempfile.TemporaryDirectory() as directory:
            docs = Path(directory, "docs")
            docs.mkdir()
            docs.joinpath("tracing.md").write_text(
                "# Tracing\n\nOpenTelemetry sends spans to Tempo for waterfall analysis.",
                encoding="utf-8",
            )
            old_db, old_docs = Config.CHECKPOINT_DB, Config.DOCUMENTATION_PATH
            Config.CHECKPOINT_DB = str(Path(directory, "checkpoints.sqlite"))
            Config.DOCUMENTATION_PATH = str(docs)
            try:
                with patch("graph.emit_log"):
                    agent = AgentGraph()
                    agent._model = FakeModel()
                    first = agent.invoke("How are spans sent to Tempo?", "thread-1")
                    second = agent.invoke("What did I ask before?", "thread-1")
                    agent.close()
            finally:
                Config.CHECKPOINT_DB, Config.DOCUMENTATION_PATH = old_db, old_docs

            self.assertEqual(first["documents"][0]["source"], "tracing.md")
            self.assertEqual(len(first["history"]), 2)
            self.assertEqual(len(second["history"]), 4)
            self.assertIn("How are spans sent to Tempo?", second["answer"])


class DocumentationStoreTest(unittest.TestCase):
    def test_returns_only_relevant_chunks(self):
        with tempfile.TemporaryDirectory() as directory:
            Path(directory, "a.md").write_text(
                "Tempo stores distributed tracing spans for waterfall inspection.", encoding="utf-8"
            )
            Path(directory, "b.md").write_text(
                "Cooking recipes describe ingredients and preparation steps.", encoding="utf-8"
            )
            results = DocumentationStore(directory).search("Tempo tracing waterfall")
            self.assertEqual([result.source for result in results], ["a.md"])


if __name__ == "__main__":
    unittest.main()

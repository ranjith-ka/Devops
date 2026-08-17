"""LangGraph workflow with documentation retrieval and durable thread memory."""

from __future__ import annotations

import os
import sqlite3
from collections.abc import Callable
from typing import Any, TypedDict

from langchain_core.messages import AIMessage, HumanMessage, SystemMessage
from langchain_ollama import ChatOllama
from langgraph.checkpoint.sqlite import SqliteSaver
from langgraph.graph import END, START, StateGraph
from opentelemetry import trace

from config import Config
from documentation import DocumentationStore
from loki import emit_log
from tracing import get_tracer


class AgentState(TypedDict, total=False):
    question: str
    documents: list[dict[str, str]]
    history: list[dict[str, str]]
    answer: str


def traced_node(name: str, function: Callable[[AgentState], dict[str, Any]]):
    """Create a graph node with a stable OpenTelemetry span and Loki events."""
    def invoke(state: AgentState) -> dict[str, Any]:
        tracer = get_tracer()
        with tracer.start_as_current_span(f"graph.node.{name}") as span:
            span.set_attribute("langgraph.node", name)
            emit_log("node.started", node=name)
            try:
                result = function(state)
                emit_log("node.completed", node=name)
                return result
            except Exception as error:
                span.record_exception(error)
                emit_log("node.failed", level="error", node=name, error_type=type(error).__name__)
                raise
    return invoke


class AgentGraph:
    """Compiled graph and its long-lived SQLite checkpointer connection."""

    def __init__(self):
        checkpoint_path = Config.CHECKPOINT_DB
        checkpoint_dir = os.path.dirname(checkpoint_path)
        if checkpoint_dir:
            os.makedirs(checkpoint_dir, exist_ok=True)
        self._connection = sqlite3.connect(checkpoint_path, check_same_thread=False)
        self._checkpointer = SqliteSaver(self._connection)
        self._documents = DocumentationStore(Config.DOCUMENTATION_PATH, Config.DOCUMENT_TOP_K)
        self._model = ChatOllama(
            model=os.getenv("OLLAMA_MODEL", "llama3.2"),
            base_url=os.getenv("OLLAMA_BASE_URL", "http://localhost:11434"),
            temperature=0,
        )
        self.graph = self._build()

    def _build(self):
        builder = StateGraph(AgentState)
        builder.add_node("retrieve_documentation", traced_node("retrieve_documentation", self._retrieve))
        builder.add_node("generate_answer", traced_node("generate_answer", self._generate))
        builder.add_node("persist_memory", traced_node("persist_memory", self._persist))
        builder.add_edge(START, "retrieve_documentation")
        builder.add_edge("retrieve_documentation", "generate_answer")
        builder.add_edge("generate_answer", "persist_memory")
        builder.add_edge("persist_memory", END)
        return builder.compile(checkpointer=self._checkpointer)

    def _retrieve(self, state: AgentState) -> dict[str, Any]:
        chunks = self._documents.search(state["question"])
        documents = [{"source": chunk.source, "text": chunk.text} for chunk in chunks]
        emit_log("documentation.retrieved", document_count=len(documents), sources=[d["source"] for d in documents])
        return {"documents": documents}

    def _generate(self, state: AgentState) -> dict[str, Any]:
        history = state.get("history", [])[-6:]
        context = "\n\n".join(
            f"Source: {document['source']}\n{document['text']}" for document in state.get("documents", [])
        ) or "No relevant local documentation was found."
        history_text = "\n".join(f"{item['role']}: {item['content']}" for item in history)
        messages = [
            SystemMessage(content=(
                "You are a concise LangChain and observability assistant. The conversation history "
                "below is durable memory restored for this thread. Use it for recall questions and "
                "never claim that you lack memory when the answer is present there. Use the supplied "
                "local documentation when relevant and mention source filenames."
            )),
            HumanMessage(content=(
                f"Conversation history:\n{history_text or '(new thread)'}\n\n"
                f"Documentation:\n{context}\n\nCurrent question: {state['question']}"
            )),
        ]
        response = self._model.invoke(messages)
        answer = response.content if isinstance(response, AIMessage) else str(response)
        emit_log("model.completed", model=os.getenv("OLLAMA_MODEL", "llama3.2"))
        return {"answer": answer}

    @staticmethod
    def _persist(state: AgentState) -> dict[str, Any]:
        history = list(state.get("history", []))
        history.extend([
            {"role": "user", "content": state["question"]},
            {"role": "assistant", "content": state["answer"]},
        ])
        emit_log("memory.persisted", message_count=len(history))
        return {"history": history[-20:]}

    def invoke(self, question: str, thread_id: str) -> AgentState:
        config = {"configurable": {"thread_id": thread_id}}
        return self.graph.invoke({"question": question}, config=config)

    def close(self) -> None:
        self._connection.close()

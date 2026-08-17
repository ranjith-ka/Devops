"""LangChain model and pipeline building for a real local Ollama model."""

from __future__ import annotations

import os
from typing import Any

from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from opentelemetry import trace

from config import Config
from tracing import get_tracer, traced_step


def build_chain():
    """Build the LangChain pipeline against a local Ollama model."""
    try:
        from langchain_ollama import ChatOllama
    except ImportError as exc:  # pragma: no cover - runtime environment check
        raise RuntimeError(
            "langchain-ollama is not installed. Run: python3 -m pip install langchain-ollama"
        ) from exc

    model_name = os.getenv("OLLAMA_MODEL", "llama3.2")
    base_url = os.getenv("OLLAMA_BASE_URL", "http://localhost:11434")

    llm = ChatOllama(
        model=model_name,
        base_url=base_url,
        temperature=0,
    )

    prompt = ChatPromptTemplate.from_messages(
        [
            ("system", "You teach LangChain in short, concrete steps."),
            ("human", "{question}"),
        ]
    )

    return (
        traced_step("prompt.render", "prompt", prompt.invoke)
        | traced_step("model.generate", "chat", llm.invoke)
        | traced_step("output.parse", "parse", StrOutputParser().invoke)
    )

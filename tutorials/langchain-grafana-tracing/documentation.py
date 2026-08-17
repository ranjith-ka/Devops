"""Small local documentation index used by the LangGraph retrieval node."""

from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class DocumentChunk:
    source: str
    text: str


class DocumentationStore:
    """Load Markdown files and retrieve relevant chunks with lexical scoring."""

    def __init__(self, path: str, top_k: int = 3):
        self.path = Path(path)
        self.top_k = top_k
        self._chunks: list[DocumentChunk] | None = None

    def _load(self) -> list[DocumentChunk]:
        chunks: list[DocumentChunk] = []
        if not self.path.exists():
            return chunks
        for file_path in sorted(self.path.rglob("*.md")):
            text = file_path.read_text(encoding="utf-8", errors="replace")
            for section in re.split(r"\n\s*\n", text):
                section = section.strip()
                if len(section) >= 40:
                    chunks.append(DocumentChunk(str(file_path.relative_to(self.path)), section[:2000]))
        return chunks

    @staticmethod
    def _terms(text: str) -> set[str]:
        return {term for term in re.findall(r"[a-z0-9_]{3,}", text.lower())}

    def search(self, query: str) -> list[DocumentChunk]:
        if self._chunks is None:
            self._chunks = self._load()
        query_terms = self._terms(query)
        if not query_terms:
            return []
        ranked = sorted(
            self._chunks,
            key=lambda chunk: len(query_terms & self._terms(chunk.text)),
            reverse=True,
        )
        return [chunk for chunk in ranked if query_terms & self._terms(chunk.text)][: self.top_k]

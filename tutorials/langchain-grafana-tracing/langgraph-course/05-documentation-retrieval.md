# Lesson 5: Documentation retrieval

## Learning objective

Understand the current local retriever and how to replace it without rewriting
the graph.

## 1. Index construction

`DocumentationStore` recursively loads `*.md` files, splits on blank lines, drops
very short sections, and caps each chunk at 2,000 characters. The index is built
lazily on the first search and cached in memory.

## 2. Ranking

The current implementation tokenizes query and chunk text, then ranks chunks by
the count of overlapping terms. It is simple, deterministic, local, and easy to
test.

It does not understand synonyms or semantic similarity. “checkpoint durability”
may not find text that only says “persist workflow state.”

## 3. Stable boundary

The graph depends only on results containing:

```python
{"source": "file.md", "text": "relevant content"}
```

That boundary allows replacement with embeddings, BM25, a vector database, or a
hybrid retriever without changing the model or memory nodes.

## Exercise

Create `docs/custom-node.md` describing a fictional `risk_check` node. Ask the UI
what `risk_check` does and confirm the source appears below the answer.

## Safe vector-store migration

1. Preserve source, section, version, and access-control metadata.
2. Filter by tenant/permissions before returning chunks.
3. Cap chunk count and total characters.
4. Trace retrieval duration and result count, not document bodies.
5. Evaluate retrieval quality with a fixed question set.

## Next lesson

Lesson 6 follows the OpenTelemetry context from the root request through every
LangGraph node.

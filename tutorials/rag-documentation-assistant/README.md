# Build a Production-Ready RAG Assistant for API Documentation

This tutorial is based on the video [Why Most RAG Systems Fail in Production | Build a Documentation Assistant for API docs](https://www.youtube.com/watch?v=--EIpKnedMc) by Applied with AI - Tamil.

The goal is to build a documentation assistant that:

- preserves the structure of technical documentation;
- finds results using both keywords and meaning;
- returns verifiable citations;
- refuses to invent an answer when the documentation is insufficient.

## Video chapters

| Time | Topic |
| --- | --- |
| [00:00](https://www.youtube.com/watch?v=--EIpKnedMc&t=0s) | Why production RAG is different |
| [02:19](https://www.youtube.com/watch?v=--EIpKnedMc&t=139s) | Document ingestion |
| [03:01](https://www.youtube.com/watch?v=--EIpKnedMc&t=181s) | A failing RAG demonstration |
| [05:00](https://www.youtube.com/watch?v=--EIpKnedMc&t=300s) | Why character chunking breaks |
| [05:46](https://www.youtube.com/watch?v=--EIpKnedMc&t=346s) | Structure-based chunking |
| [06:51](https://www.youtube.com/watch?v=--EIpKnedMc&t=411s) | RAG application flow |
| [07:30](https://www.youtube.com/watch?v=--EIpKnedMc&t=450s) | Citation metadata design |
| [11:05](https://www.youtube.com/watch?v=--EIpKnedMc&t=665s) | Code walkthrough |
| [16:07](https://www.youtube.com/watch?v=--EIpKnedMc&t=967s) | Final demonstration |

## 1. Why a simple RAG pipeline fails

A basic RAG tutorial often follows this sequence:

1. Extract text from a document.
2. Split the text every fixed number of characters.
3. Create embeddings.
4. Retrieve the nearest chunks.
5. Ask an LLM to answer.

That approach can work for prose, but API documentation has meaningful boundaries: pages, headings, endpoint names, parameter tables, examples, warnings, and code blocks. A fixed-size splitter can detach a parameter from its endpoint, split a code sample in half, or discard the heading that gives a paragraph its meaning.

It also struggles with exact identifiers. Semantic search may understand that “remove a user” is related to “delete an account,” but it may rank the wrong result for an exact token such as `DELETE /v1/users/{id}` or `rate_limit_remaining`.

The production design used in this tutorial addresses both problems:

```text
Documentation
      |
      v
Parse structure -> create metadata-rich chunks -> build two indexes
                                                |              |
                                                v              v
User question -----------------------------> BM25 search   Vector search
                                                |              |
                                                +------v-------+
                                                       RRF
                                                        |
                                                        v
                                            grounded answer + citations
```

## 2. Project layout

```text
docs-assistant/
├── data/
│   └── docs/
├── src/
│   ├── ingest.py
│   ├── retrieve.py
│   └── answer.py
├── tests/
│   └── questions.json
├── .env
└── requirements.txt
```

The examples below are intentionally framework-neutral. Replace the parser, embedding model, vector database, or LLM without changing the architecture.

## 3. Design the chunk schema first

Do not store only the chunk text. Citations and filtering depend on metadata captured during ingestion.

```python
from dataclasses import dataclass, field


@dataclass
class DocumentChunk:
    id: str
    content: str
    title: str
    section: str
    source_url: str
    document_id: str
    chunk_index: int
    content_type: str = "text"  # text, code, table, warning
    metadata: dict = field(default_factory=dict)
```

Useful metadata includes:

- stable document and chunk IDs;
- page title and heading path;
- canonical source URL, including a heading anchor when possible;
- product and API version;
- endpoint and HTTP method;
- programming language for code samples;
- ingestion timestamp or source revision.

A model cannot produce trustworthy citations if the retriever never receives trustworthy source metadata.

## 4. Parse documents by structure

For Markdown, headings provide natural semantic boundaries. Keep the heading hierarchy with every chunk.

```python
import re
from collections.abc import Iterable


HEADING = re.compile(r"^(#{1,6})\s+(.+)$")


def split_markdown_by_heading(markdown: str) -> Iterable[tuple[str, str]]:
    """Yield (heading_path, body) sections without losing context."""
    headings: list[str] = []
    body: list[str] = []

    def current_path() -> str:
        return " > ".join(headings) if headings else "Overview"

    for line in markdown.splitlines():
        match = HEADING.match(line)
        if not match:
            body.append(line)
            continue

        if body and any(part.strip() for part in body):
            yield current_path(), "\n".join(body).strip()
        body = []

        level = len(match.group(1))
        heading = match.group(2).strip()
        headings = headings[: level - 1]
        headings.append(heading)

    if body and any(part.strip() for part in body):
        yield current_path(), "\n".join(body).strip()
```

Large sections may still need subdivision. Split only inside a section, preferably at paragraph boundaries, and repeat the heading path in each child chunk.

```python
def split_large_section(section: str, max_words: int = 350) -> list[str]:
    paragraphs = section.split("\n\n")
    chunks: list[str] = []
    current: list[str] = []
    word_count = 0

    for paragraph in paragraphs:
        paragraph_words = len(paragraph.split())
        if current and word_count + paragraph_words > max_words:
            chunks.append("\n\n".join(current))
            current, word_count = [], 0
        current.append(paragraph)
        word_count += paragraph_words

    if current:
        chunks.append("\n\n".join(current))
    return chunks
```

Production parsers should additionally keep fenced code blocks and tables intact. For HTML documentation, use the document object model rather than stripping every tag into one text stream.

## 5. Create searchable text

Add structural context before embedding and keyword indexing:

```python
def searchable_text(chunk: DocumentChunk) -> str:
    return (
        f"Title: {chunk.title}\n"
        f"Section: {chunk.section}\n"
        f"Content type: {chunk.content_type}\n\n"
        f"{chunk.content}"
    )
```

This makes a paragraph under `Authentication > API keys > Rotation` searchable even if the paragraph itself contains only “Rotate them every 90 days.”

Use the same normalized chunk collection to build:

1. a lexical index for BM25;
2. a vector index for semantic similarity.

## 6. Retrieve with hybrid search

### BM25 retrieval

BM25 is strong at exact terms: class names, error codes, field names, paths, and command-line flags.

```python
def bm25_search(query: str, index, chunks: list[DocumentChunk], limit: int = 20):
    tokens = query.lower().split()
    scores = index.get_scores(tokens)
    ranked = sorted(range(len(scores)), key=lambda i: scores[i], reverse=True)
    return [(chunks[i], float(scores[i])) for i in ranked[:limit]]
```

### Semantic retrieval

Vector search is strong when the question and the documentation express the same concept with different words.

```python
def semantic_search(query: str, embed, vector_store, limit: int = 20):
    query_vector = embed(query)
    return vector_store.search(vector=query_vector, limit=limit)
```

Retrieve more candidates than the final context requires. Fusion needs a useful candidate pool; the final prompt does not need all of it.

## 7. Combine rankings with Reciprocal Rank Fusion

Raw BM25 and vector similarity scores are not directly comparable. Reciprocal Rank Fusion (RRF) combines positions instead of score magnitudes:

\[
RRF(d) = \sum_{r \in rankings} \frac{1}{k + rank_r(d)}
\]

`k = 60` is a common starting point. Tune it with your evaluation set.

```python
from collections import defaultdict


def reciprocal_rank_fusion(result_lists, k: int = 60):
    scores = defaultdict(float)
    chunks_by_id = {}

    for results in result_lists:
        for rank, result in enumerate(results, start=1):
            chunk = result[0] if isinstance(result, tuple) else result
            chunks_by_id[chunk.id] = chunk
            scores[chunk.id] += 1.0 / (k + rank)

    ranked_ids = sorted(scores, key=scores.get, reverse=True)
    return [(chunks_by_id[chunk_id], scores[chunk_id]) for chunk_id in ranked_ids]
```

End-to-end retrieval then becomes:

```python
def retrieve(query: str, bm25_index, chunks, embed, vector_store, top_k: int = 6):
    lexical = bm25_search(query, bm25_index, chunks, limit=20)
    semantic = semantic_search(query, embed, vector_store, limit=20)
    fused = reciprocal_rank_fusion([lexical, semantic])
    return fused[:top_k]
```

Before generation, optionally remove duplicate or nearly identical chunks and cap the number of chunks from any single page.

## 8. Build citation-ready context

Give every retrieved chunk a short source label. The label, title, section, and URL must be supplied by your code—not invented by the model.

```python
def build_context(results) -> str:
    blocks = []
    for number, (chunk, _score) in enumerate(results, start=1):
        blocks.append(
            f"[S{number}]\n"
            f"Title: {chunk.title}\n"
            f"Section: {chunk.section}\n"
            f"URL: {chunk.source_url}\n"
            f"Content:\n{chunk.content}"
        )
    return "\n\n---\n\n".join(blocks)
```

Render the final citation links from retrieved metadata in the application layer. Treat an LLM citation such as `[S2]` only as a reference to the supplied source map.

## 9. Enforce grounded answers

Use a strict system instruction:

```text
You are an assistant for API documentation.

Answer only from the supplied sources.
Every factual claim must cite one or more source labels such as [S1].
Do not use outside knowledge to fill gaps.
If the sources do not contain enough information, respond:
"I don't know based on the available documentation."
Do not create URLs, options, parameters, defaults, or code that are absent
from the sources.
```

Then construct the request:

```python
def make_prompt(question: str, results) -> str:
    return f"""Sources:
{build_context(results)}

Question: {question}

Return a concise answer with inline source labels.
"""
```

The refusal instruction is necessary but not sufficient. Add application-level controls:

- require at least one citation for a factual answer;
- verify that every returned source label exists in the supplied source map;
- return the refusal response when retrieval is empty or below a tuned threshold;
- log the question, retrieved IDs, ranks, answer, and citations;
- never let the model generate source URLs directly.

## 10. Example behavior

Question:

```text
How do I authenticate requests, and how often should I rotate the key?
```

Good answer:

```text
Send the API key in the `Authorization` header using the format shown in the
authentication guide [S1]. Rotate the key every 90 days [S2].
```

If rotation frequency is absent from the retrieved documentation, the assistant should not propose a general security recommendation as though it came from the API docs. It should say that the available documentation does not specify the frequency.

## 11. Evaluate before shipping

Create a small, version-controlled test set containing:

- direct questions whose answers appear in one section;
- questions requiring evidence from multiple sections;
- exact identifiers and error codes;
- paraphrased questions;
- ambiguous questions;
- unanswerable questions;
- questions about an older or unsupported API version.

Measure each stage separately.

| Layer | Useful checks |
| --- | --- |
| Ingestion | Heading, code block, table, URL, and version preservation |
| Retrieval | Recall@k, Mean Reciprocal Rank, exact-identifier coverage |
| Answer | Correctness, completeness, citation precision, citation coverage |
| Safety | Refusal accuracy and unsupported-claim rate |
| Operations | Latency, token usage, index freshness, error rate |

When an answer is wrong, first inspect the retrieved chunks. Prompt changes cannot repair missing or badly split evidence.

## 12. Production hardening checklist

- [ ] Crawl only canonical, approved documentation sources.
- [ ] Preserve heading hierarchy, code blocks, tables, and warnings.
- [ ] Record product and API versions in chunk metadata.
- [ ] Use stable IDs so changed pages can be re-indexed safely.
- [ ] Delete stale chunks when a source document is removed.
- [ ] Combine lexical and semantic retrieval.
- [ ] Tune retrieval and refusal thresholds using labeled questions.
- [ ] Map citations to URLs outside the LLM.
- [ ] Reject unknown citation labels.
- [ ] Keep tenant or access-control filters in the retrieval layer.
- [ ] Protect the prompt from instructions found inside retrieved documents.
- [ ] Log retrieval traces without leaking credentials or private content.
- [ ] Re-run evaluation after documentation, embedding, prompt, or model changes.

## 13. Common failure modes

### The correct page is indexed but never retrieved

Check whether its title and heading path were included in the searchable text. Confirm that the query's exact identifiers reach BM25 without destructive tokenization.

### The right chunk appears below irrelevant chunks

Inspect the lexical and semantic rankings independently. Adjust candidate counts, metadata filters, RRF settings, or add a reranking stage.

### Citations link to a page that does not support the claim

Use smaller structure-aware chunks, require claim-level citations, and measure citation correctness separately from answer fluency.

### The assistant answers unsupported questions

Add negative examples to the evaluation set, strengthen the refusal contract, tune the retrieval confidence gate, and validate citations after generation.

### Newly published documentation is missing

Track the source revision and ingestion time. Run incremental ingestion and expose index freshness in monitoring.

## Key takeaway

Production RAG is primarily an information-retrieval and data-quality system. The LLM is the final presentation layer. Reliable results come from preserving document structure, combining complementary retrieval methods, carrying source metadata through the entire pipeline, and making “I don’t know” a valid outcome.

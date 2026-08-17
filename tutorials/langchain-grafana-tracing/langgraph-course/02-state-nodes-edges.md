# Lesson 2: Typed state, nodes, and edges

## Learning objective

Understand how `StateGraph` turns small Python functions into an explicit,
testable workflow.

## 1. The state contract

`AgentState` in `graph.py` is the shared contract:

```python
class AgentState(TypedDict, total=False):
    question: str
    documents: list[dict[str, str]]
    history: list[dict[str, str]]
    answer: str
```

Each node reads existing fields and returns only the fields it updates. For
example, documentation retrieval returns `{"documents": documents}`. LangGraph
merges that partial update into the state.

`total=False` allows intermediate states to omit fields that have not been
produced yet. It does not make arbitrary fields safe; every durable value should
still be declared.

## 2. Node responsibilities

The current graph follows single-responsibility boundaries:

- `retrieve_documentation`: question in, document chunks out;
- `generate_answer`: question + documents + history in, answer out;
- `persist_memory`: question + answer + old history in, new history out.

Small nodes are easier to retry, trace, test, replace, and compare.

## 3. Edges define control flow

```python
builder.add_edge(START, "retrieve_documentation")
builder.add_edge("retrieve_documentation", "generate_answer")
builder.add_edge("generate_answer", "persist_memory")
builder.add_edge("persist_memory", END)
```

This is a fixed graph. The model cannot skip retrieval or memory persistence.
That is appropriate while the workflow is predictable.

## Exercise

Add a state field named `request_category`. Create a temporary node that returns
`{"request_category": "documentation"}` and place it before retrieval. Confirm
that a new `graph.node.classify_request` span appears in Tempo.

## Common mistakes

- Mutating the incoming state instead of returning a partial update.
- Letting one node retrieve, call the model, persist memory, and format output.
- Adding a state value without updating `AgentState`.
- Calling external systems without timeouts.

## Next lesson

Lesson 3 connects browser requests to graph invocations and explains why thread
identity must be separate from trace identity.

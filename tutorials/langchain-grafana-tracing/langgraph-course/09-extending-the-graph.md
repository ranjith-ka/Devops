# Lesson 9: Extending the graph safely

## Learning objective

Add new behavior while preserving state contracts, trace structure, logging,
memory, and tests.

## Recipe 1: Add a deterministic node

1. Add its output field to `AgentState`.
2. Write a function that accepts state and returns a partial update.
3. Register it through `traced_node`.
4. Add edges.
5. Test the state update.
6. Confirm its Tempo span and Loki events.

Example validation node:

```python
def _validate(self, state: AgentState) -> dict[str, Any]:
    answer = state.get("answer", "").strip()
    if not answer:
        raise ValueError("Model returned an empty answer")
    return {"answer": answer}
```

Place it between generation and memory persistence so invalid answers are not
stored.

## Recipe 2: Add conditional routing

Use a small router that returns a bounded route value, then conditional edges:

```text
classify_request
  |- documentation -> retrieve_documentation
  `- general       -> generate_answer
```

Keep routing decisions in state and logs. Do not let arbitrary model text become
a node name; map allowed decisions explicitly.

## Recipe 3: Add tools

Wrap each external call in its own node or tool boundary. Validate arguments,
apply authorization, set timeouts, bound retries, and return structured data.
Never place credentials in graph state, spans, or logs.

## Recipe 4: Add retries

Retry only transient operations. Record attempt number and final status. Repeated
spans will appear as occurrence keys in comparison reports.

## Suggested next nodes

1. `classify_request`
2. `validate_answer`
3. `retrieve_vector_documents`
4. `execute_tool`
5. `summarize_history`

Add one at a time and preserve the tests after each change.

## Exercise

Implement `validate_answer`, add a test for an empty answer, and verify that a
failed validation creates an error span plus `node.failed` Loki event.

## Next lesson

Lesson 10 turns this single-process tutorial into a production roadmap.

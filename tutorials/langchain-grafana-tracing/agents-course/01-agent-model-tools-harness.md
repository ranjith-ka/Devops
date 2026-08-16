# Lesson 1: Agent, model, tools, and harness

## Learning objective

By the end of this lesson, you should be able to explain:

- how an agent differs from a single model call;
- what the model and tools each do;
- what the agent loop does;
- why LangChain calls the surrounding configuration a harness;
- when a fixed chain is a better choice than an agent.

## 1. Model call versus agent

A single model call follows a direct path:

```text
User input -> model -> answer
```

An agent can repeat a model-and-tool loop:

```text
User request
     |
     v
   Model <------------------+
     |                      |
     v                      |
Does it need a tool?        |
  |                         |
  +-- no --> final answer   |
  |                         |
  `-- yes --> call tool ----+
```

The model chooses whether a tool is needed. LangChain executes the selected tool,
adds its result to the agent state, and calls the model again. The loop ends when
the model returns a final response instead of another tool call.

## 2. The four foundational pieces

### Model

The model is the reasoning component. It reads messages and decides whether to
answer or request a tool call.

LangChain commonly identifies models with a `provider:model` string:

```python
model = "openai:YOUR_MODEL_NAME"
```

You can also pass an initialized model object. We will compare those approaches in
Lesson 3.

### Tools

Tools are functions the agent is allowed to call:

```python
from langchain.tools import tool


@tool
def get_temperature(city: str) -> str:
    """Return the current temperature for a city."""
    return f"The temperature in {city} is 28 C."
```

The function name, type hints, and docstring help the model understand when and
how to call the tool. The model does not receive permission to execute arbitrary
Python merely because it is called an agent.

### Agent loop

For the request `What is the temperature in Chennai?`, a run might be:

1. The model receives the user message.
2. The model requests `get_temperature(city="Chennai")`.
3. LangChain validates and executes the tool call.
4. LangChain appends the tool result to the messages.
5. The model sees the updated messages.
6. The model produces the final answer.

`create_agent` constructs this loop, so application code does not need to manually
alternate between model and tool calls.

### Harness

The harness is everything that shapes and supports the loop:

```text
Harness
|- model
|- tools
|- system prompt
|- state and memory
|- middleware
|- guardrails
|- retry policies
`- execution environment
```

A useful mental model is:

```text
Agent = model + harness
```

The model supplies general reasoning ability. The harness determines what the
agent knows, what it can do, what it must not do, and how failures are handled.

## 3. Preview: a minimal agent

Do not worry about running this example yet. Installation, credentials, and model
selection are covered in Lesson 2.

```python
from langchain.agents import create_agent
from langchain.tools import tool


@tool
def get_temperature(city: str) -> str:
    """Return the current temperature for a city."""
    return f"The temperature in {city} is 28 C."


agent = create_agent(
    model="openai:YOUR_MODEL_NAME",
    tools=[get_temperature],
    system_prompt="You are a concise weather assistant.",
)

result = agent.invoke(
    {
        "messages": [
            {
                "role": "user",
                "content": "What is the temperature in Chennai?",
            }
        ]
    }
)

print(result["messages"][-1].content)
```

The three important inputs to `create_agent` are already visible:

- `model`: the reasoning engine;
- `tools`: the allowed actions;
- `system_prompt`: instructions that shape behavior.

## 4. Agent versus chain

A chain follows a path selected by the developer:

```text
prompt -> model -> parser
```

An agent dynamically selects its next action:

```text
model -> perhaps a tool -> model -> perhaps another tool -> answer
```

Use a chain when the workflow is predictable. Use an agent when the correct next
step depends on information discovered during execution.

| Task | Better choice | Reason |
| --- | --- | --- |
| Translate text | Chain | One predictable operation |
| Extract an invoice into JSON | Chain | Fixed schema and path |
| Diagnose a server incident | Agent | Each observation changes the next step |
| Research and compare products | Agent | Search and comparison steps vary |

Agents add model calls, latency, cost, and nondeterminism. Do not use one when a
small deterministic pipeline solves the problem.

## 5. Common mistakes

### Calling every chatbot an agent

A model with only a system prompt is still essentially a chatbot. The important
agent behavior is selecting actions and continuing the loop based on results.

### Giving tools vague descriptions

The model chooses tools using their names, schemas, and descriptions. A vague
docstring produces unreliable selection.

### Giving an agent excessive authority

Only expose actions the use case needs. Read operations, reversible operations,
and destructive operations should have different controls.

### Using an agent for a fixed workflow

If every request must perform the same three steps, encode those steps as a chain
or graph. It will usually be cheaper and easier to test.

## 6. Exercise

Choose `chain` or `agent` for each task before revealing the answers:

1. Translate English into Tamil.
2. Check several websites and compare prices.
3. Convert an invoice into a fixed JSON schema.
4. Investigate a failed Kubernetes deployment using logs and cluster tools.
5. Classify a support ticket into one of five categories.

<details>
<summary>Suggested answers</summary>

1. **Chain**: the operation is predictable.
2. **Agent**: it may choose different searches and comparison steps.
3. **Chain**: use structured output with a fixed schema.
4. **Agent**: each diagnostic result determines the next action.
5. **Chain**: classification is normally a single predictable model operation.

</details>

## 7. Knowledge check

You are ready for Lesson 2 if you can answer these without looking back:

1. What causes the agent loop to stop?
2. Who selects a tool: LangChain application code or the model?
3. What does LangChain do after a tool returns?
4. Name four things that belong to the harness.
5. Why can a chain be better than an agent?

## Next lesson

Lesson 2 will install the required packages, configure a model provider safely,
create an agent with `create_agent`, invoke it, and inspect the returned messages.


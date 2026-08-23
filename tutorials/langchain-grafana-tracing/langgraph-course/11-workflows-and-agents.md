# 11. Workflows, agents, and the next incident-investigation graph

This lesson connects the common LangGraph patterns to the code in this repository.
Read it after [Lesson 10](./10-production-roadmap.md).

Reference: [LangGraph workflows and agents](https://docs.langchain.com/oss/python/langgraph/workflows-agents).

## 1. Start with the mental model

A LangGraph application is built from four concepts:

- **State** carries information through a run.
- **Nodes** perform work and return partial state updates.
- **Edges** decide which node runs next.
- **START** and **END** are virtual entry and exit markers; they do not execute
  application logic.

Calling `builder.add_node(...)` registers a function. Calling `compile()`
validates and prepares the graph. Nodes execute only after `graph.invoke(...)`
or a streaming method starts a run.

The current implementation in [`graph.py`](../graph.py) is:

```text
START -> retrieve_documentation -> generate_answer -> persist_memory -> END
```

Its `AgentState` contains the question, retrieved documents, history, and answer.
Despite the class name `AgentGraph`, this is currently a **deterministic
workflow**: application code controls the complete route.

## 2. Workflow versus agent

A workflow follows paths selected by application code. An agent lets the model
choose actions from a bounded toolset and repeat the reason/action/observation
loop until it can answer.

```text
Workflow: START -> retrieve -> generate -> persist -> END

Agent:    START -> model -> tool? -> model -> tool? -> model -> END
                              ^                    |
                              +--------------------+
```

An LLM call inside a node does not, by itself, make the graph an agent. The
important question is: **who selects the next action at runtime?**

## 3. Workflow patterns

| Pattern | Meaning | Possible tracing use case |
| --- | --- | --- |
| Prompt chaining | One LLM result feeds the next fixed step | analyse -> review -> final RCA |
| Parallelization | Independent predefined nodes run together | latency, errors, and database analysis |
| Routing | Select one specialized branch | send HTTP 500 to error analysis |
| Orchestrator-worker | Create an unknown number of tasks at runtime | one worker per suspicious service |
| Evaluator-optimizer | Evaluate and improve until accepted or limited | require evidence and next actions in RCA |

Routing is normally implemented with conditional edges. Conditional edges are
the mechanism; routing is the design pattern.

Parallel nodes and worker tasks are not automatically agents. They become more
agentic only when an LLM dynamically decides the task decomposition or actions.

## 4. Tool-calling agent loop

For this POC, the useful next graph is an incident-investigation agent:

```text
User question
    |
    v
Model bound to approved tools
    |
    +-- tool call? -- yes --> ToolNode --> tool result --> model
    |
    +-- no ---------------------------------------------> END
```

Recommended read-only tools:

- `get_trace(trace_id)` backed by `fetch_trace_from_tempo`;
- `get_logs(trace_id)` backed by `fetch_logs_for_trace`;
- `compare_traces(baseline_trace_id, candidate_trace_id)`;
- later, `get_deployment(service, environment)` backed by Kubernetes or Flux.

`bind_tools(tools)` advertises tool schemas to the model. It does not execute
them. `ToolNode(tools)` executes requested tools and returns `ToolMessage`
objects. `tools_condition` routes to the tool node when the latest model
message contains tool calls; otherwise it routes to `END`.

## 5. State, runtime context, and secrets

Use the correct source for every value:

| Value | Source |
| --- | --- |
| Trace ID selected from the question | Model-generated tool argument |
| Conversation and investigation messages | `MessagesState` |
| Environment and cluster for the run | Custom graph state |
| Grafana tenant or organization | Runtime context |
| Grafana token | Workload identity, environment, or secret manager |

Do not expose credentials as tool arguments or return them in tool results.
Arguments are model-controlled and therefore untrusted.

## 6. Safety boundaries

The first agent POC should remain read-only. Add:

- a recursion or iteration limit;
- HTTP timeouts and bounded result sizes;
- strict trace-ID validation;
- tool error messages that do not expose credentials;
- evidence citations naming the trace, span, log event, or deployment;
- a rule that the model reports uncertainty instead of inventing a root cause;
- human approval before any future write action.

## 7. Completion exercise

Implement the milestone described in
[`POC-COMPLETION.md`](../POC-COMPLETION.md). The first successful demonstration
should prove that the model chooses tools dynamically rather than following a
hard-coded `Tempo -> Loki -> answer` sequence.

You should be able to show two runs:

1. A healthy trace where the agent calls only `get_trace` and stops.
2. A failing trace where the agent calls `get_trace`, then `get_logs`, and
   produces an evidence-backed root-cause hypothesis.

That difference in runtime paths is the simplest proof that the POC contains a
real agent loop.

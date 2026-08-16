# LangChain Agents: step-by-step course

This course follows the concepts in the current
[LangChain Agents documentation](https://docs.langchain.com/oss/python/langchain/agents).
Each lesson contains a plain-language explanation, mental model, example,
mistakes to avoid, and an exercise.

## How to use this course

1. Read one lesson at a time.
2. Type and run its examples instead of copying them blindly.
3. Complete the exercise before moving forward.
4. Add tracing only after the basic agent behavior is clear.

## Curriculum

| Lesson | Concept | Status |
| --- | --- | --- |
| [01](./01-agent-model-tools-harness.md) | Agent, model, tools, and harness | Ready |
| 02 | Installation and `create_agent` | Planned |
| 03 | Models and model providers | Planned |
| 04 | Tools and the `@tool` decorator | Planned |
| 05 | System prompts | Planned |
| 06 | The agent execution loop | Planned |
| 07 | Invocation and messages | Planned |
| 08 | Conversation history and checkpointers | Planned |
| 09 | Agent state | Planned |
| 10 | Runtime context | Planned |
| 11 | Structured output | Planned |
| 12 | Streaming | Planned |
| 13 | Middleware fundamentals | Planned |
| 14 | Execution environments and filesystem tools | Planned |
| 15 | Context management and summarization | Planned |
| 16 | Memory and skills | Planned |
| 17 | Planning with todo lists | Planned |
| 18 | Delegation and subagents | Planned |
| 19 | Naming agents | Planned |
| 20 | Fault tolerance and retries | Planned |
| 21 | Guardrails and PII protection | Planned |
| 22 | Steering agent behavior | Planned |
| 23 | Observability and Grafana tracing | Planned |
| 24 | Scaling and cost optimization | Planned |

## Course project

Across the lessons, we will grow a small support agent:

```text
User request
     |
     v
Support agent
  |- knowledge search tool
  |- service-status tool
  |- structured response
  |- conversation memory
  |- retries and guardrails
  `- OpenTelemetry traces -> Grafana Tempo
```

The early lessons use deterministic local functions where possible. Provider API
calls are introduced separately so you can distinguish LangChain concepts from
provider configuration and cost.


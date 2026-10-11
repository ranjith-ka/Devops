# AI DevOps Team

Planning scaffold for a multi-agent DevOps workflow using Codex CLI.

Start with [project context](docs/chatgpt-context.md), then read the
[implementation plan](docs/plan.md). No runnable agent configuration or
infrastructure has been created yet.

## Intended workflow

Request → coordinator plan → specialist tasks → integration → QA and security
review → documented outcome → human approval for production → deployment verification.

## Contents

- `AGENTS.md`: project working instructions.
- `docs/chatgpt-context.md`: known requirements and missing environment details.
- `docs/plan.md`: phases, responsibilities, validation, and approval gates.
- `templates/task.md`: delegation and handoff template.
- `templates/outcome.md`: validation and decision record.

## First session

From this directory, start your installed `codex` CLI and use:

> Read AGENTS.md and docs/chatgpt-context.md. Complete phase 0 of docs/plan.md
> using read-only discovery. Delegate independent inventory tasks. Report
> missing facts and propose the first local pilot before changing agent
> configuration or infrastructure.

Codex supports specialized subagents and instructions requesting delegation.
Confirm capabilities against the installed CLI before configuring roles.
See [official subagent documentation](https://learn.chatgpt.com/docs/agent-configuration/subagents).

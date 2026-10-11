# Project context

## Provenance

The workspace's requested `docs/chatgpt-context.md` was absent when this scaffold
was created. This document records only the user's supplied requirements; it
does not reconstruct the inaccessible earlier ChatGPT conversation.

## Goal

Build a multi-agent DevOps team using Codex CLI, with these specialists:

- Kubernetes and GitOps
- CI/CD
- SRE and Observability
- Security
- Automation
- QA

## Working principles

Plan before implementation. Delegate independent tasks. Require approval for
production changes. Validate changes with tests. Document decisions and outcome.

## Facts to establish before architecture or configuration

- Installed Codex CLI version, supported subagent configuration, and authentication.
- Target repository and application, languages, and existing test commands.
- Cloud provider, Kubernetes distributions, cluster contexts, and namespaces.
- GitOps controller and repository; CI/CD provider and deployment triggers.
- Observability stack, SLOs, alert routing, and incident workflow.
- Secret management, identity, and permitted tools for each environment.
- Definition of production, approver, approval mechanism, and rollback ownership.
- Task concurrency, token/cost limits, and first pilot acceptance criteria.

## Status

Planning only. Proposed workflow in `plan.md` is provisional pending discovery.
No production access, deployment, provider selection, or agent runtime setup is
implied by this scaffold.

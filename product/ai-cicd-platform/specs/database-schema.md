# Database Schema

## Core Entities

- organizations
- users
- memberships
- repositories
- pipelines
- pipeline_runs
- pipeline_steps
- deployments
- deployment_runs
- artifacts
- logs
- metrics_snapshots
- pull_requests
- security_findings
- cost_events
- ai_threads
- ai_messages
- ai_memories
- prompt_versions
- tool_calls
- notifications
- audit_events

## Relationship Summary

- Organization owns repositories, users, and policies.
- Repository owns pipelines and pull requests.
- Pipeline owns runs and steps.
- Deployment owns deployment runs and rollback actions.
- AI thread is tenant-scoped and linked to evidence references.
- Prompt versions and tool calls are immutable audit artifacts.

## Storage Decisions

- PostgreSQL for tenancy, metadata, state, and audit trails.
- Redis for hot cache, summaries, sessions, and queue coordination.
- Object storage for logs, artifacts, screenshots, and generated files.
- OpenSearch for text search over logs, PRs, and release notes.
- Qdrant for embeddings and semantic retrieval.

## Suggested Table Notes

- `organizations`: id, name, slug, plan, created_at.
- `users`: id, email, display_name, status, created_at.
- `memberships`: org_id, user_id, role, permissions, invited_at.
- `repositories`: org_id, provider, name, default_branch, sync_state.
- `pipelines`: repository_id, provider, name, yaml_path, last_synced_at.
- `pipeline_runs`: pipeline_id, commit_sha, status, started_at, finished_at, duration_ms.
- `deployments`: repository_id, environment, version, status, risk_score.
- `ai_threads`: org_id, subject_type, subject_id, created_by, memory_summary.
- `ai_messages`: thread_id, role, content, citations, created_at.
- `audit_events`: org_id, actor_id, action, resource_type, resource_id, metadata.

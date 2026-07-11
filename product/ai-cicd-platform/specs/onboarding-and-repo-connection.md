# Next Phase: User Onboarding and Repository Connection

## Goal

Convert first-time visitors into connected organizations with real delivery data flowing into the platform.

## Product Outcome

A new user can sign in, create or join an organization, connect a repository provider, authorize access, select a repository, and land on a populated dashboard with AI-ready data sources.

## Primary User Journey

1. User opens the landing page.
2. User clicks `Get started` or `Connect repository`.
3. User signs in or creates an organization.
4. User chooses a source provider such as GitHub or Azure DevOps.
5. User authorizes the app and selects one or more repositories.
6. User confirms the first sync scope and sees progress.
7. User lands on the dashboard with repo health, pipeline summary, and AI guidance.

## Onboarding Screens

### 1. Welcome Screen

Purpose:
- Explain what the product does.
- Show the value of connecting a repo.
- Offer sign-in and trial entry points.

States:
- Loading: skeleton hero and disabled CTAs.
- Error: auth provider unavailable or misconfigured.
- Empty: no connected workspace yet.

### 2. Organization Setup

Purpose:
- Create organization name, slug, and default region.
- Pick the initial plan or trial mode.
- Confirm team size and primary use case.

States:
- Loading: form skeleton.
- Error: validation and uniqueness failures.
- Empty: no org selected yet.

### 3. Repository Connection Wizard

Purpose:
- Select provider.
- Authorize access.
- Choose repositories.
- Set sync scope and initial import depth.

States:
- Loading: provider handshake and repo list skeleton.
- Error: auth denied, scope denied, or API rate limit.
- Empty: no repositories found.

### 4. Sync Progress

Purpose:
- Show ingestion progress for commits, pipelines, PRs, logs, and deployment history.
- Surface what data is ready and what is still pending.

States:
- Loading: progress bars, live events, and estimated time.
- Error: partial sync or provider timeout.
- Empty: nothing imported yet.

### 5. First Value Dashboard

Purpose:
- Show repository health, pipeline summary, and first AI-generated insight.
- Push the user into an actionable next step.

States:
- Loading: cached summary and charts.
- Error: backfill or retrieval failure.
- Empty: connected but not yet enough data for a meaningful summary.

## Core Functional Requirements

- Support onboarding from landing page or invite link.
- Support GitHub and Azure DevOps as initial repository providers.
- Support organization-level and repository-level permission scopes.
- Support single repo and multi repo connection.
- Support dry-run validation before full sync.
- Support retry after partial failure.
- Save onboarding state so the flow can resume after refresh or logout.

## API and Contract Needs

### Auth and Onboarding

- `POST /auth/start`
- `POST /auth/callback`
- `POST /organizations`
- `GET /organizations/:id/onboarding-state`
- `PATCH /organizations/:id/onboarding-state`

### Repository Connection

- `POST /integrations/github/connect`
- `POST /integrations/azure-devops/connect`
- `GET /integrations/providers`
- `GET /integrations/repositories`
- `POST /integrations/repositories/sync`

### Suggested Response Fields

- `status`
- `step`
- `progress`
- `error_code`
- `error_message`
- `connected_repositories`
- `next_action`

## Data Model

New or extended entities:
- `organization_onboarding`
- `integration_providers`
- `integration_connections`
- `integration_repositories`
- `repository_sync_jobs`
- `repository_sync_events`

Important fields:
- provider name
- tenant/org id
- installation id or OAuth token reference
- sync scope
- last sync status
- first sync completed at

## Event Model

New events:
- `onboarding.started`
- `onboarding.completed`
- `integration.connected`
- `integration.authorized`
- `repository.selected`
- `repository.sync.started`
- `repository.sync.completed`
- `repository.sync.failed`

## UX Requirements

- Keep the onboarding flow linear and short.
- Allow users to skip advanced config and return later.
- Make the repo connection wizard feel safe and explicit about permissions.
- Show what data will be ingested before the user confirms.
- Include a visible resume state when the user leaves midway.

## Acceptance Criteria

- A new user can reach a connected repository without needing support.
- The platform clearly explains what is imported and why.
- The first dashboard load contains real repo and pipeline data.
- The onboarding state survives refresh and browser restart.
- Partial sync failures are visible and recoverable.

## Implementation Order

1. Auth entry points and organization creation.
2. Integration provider discovery and repository selection.
3. Sync job orchestration and progress reporting.
4. First-value dashboard and AI summary generation.
5. Resumable onboarding state and analytics.

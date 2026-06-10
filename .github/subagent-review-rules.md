# Subagent Usage Rules

## Purpose
Use subagents to parallelize read-heavy tasks (search, diff analysis, review) and keep the main agent focused on decisions and final output quality.

## What Agents Can Do
1. Explore code and docs quickly across the repository.
2. Analyze diffs and summarize behavioral impact.
3. Perform security, correctness, and reliability checks.
4. Validate assumptions with file/line evidence.
5. Propose focused fixes with minimal change scope.

## What Agents Must Not Do
1. Merge, rebase, force-push, or rewrite history.
2. Run destructive commands (for example hard reset, branch deletion) without explicit approval.
3. Expose secrets from env files, tokens, or credentials in outputs.
4. Claim checks passed without running or validating evidence.
5. Ignore base branch context when doing PR review.

## When To Call Subagents
1. Large change sets (many files or mixed domains).
2. Need independent verification of risky logic.
3. Need fast repository-wide search before editing.
4. Need a dedicated review pass before PR submission.

## PR Review Workflow (Subagent-First)
1. Determine review range: compare current branch against base branch (main by default).
2. Gather scope: commit list, changed files, and diff stats.
3. Run subagent review with explicit rubric:
- correctness/regressions
- security/privacy
- performance/reliability
- docs/tests gaps
4. Return findings ordered by severity:
- Critical
- High
- Medium
- Low
5. Each finding should include:
- file path
- line reference
- why it matters
- concrete fix suggestion
6. If no findings: state that explicitly and list residual risks/testing gaps.

## Review Rubric
1. Behavior changed unintentionally?
2. Inputs validated and errors handled?
3. Security controls present (auth, secrets, transport)?
4. Timeouts/retries/resource cleanup handled?
5. Tests and docs updated for the change?

## Output Contract
1. Findings first, summary second.
2. No vague statements: include evidence.
3. Prefer minimal, actionable fixes.
4. Call out assumptions and unknowns clearly.

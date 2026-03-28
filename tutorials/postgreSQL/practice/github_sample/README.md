# GitHub-style SQL practice data

Synthetic **PostgreSQL** tables modeled after repositories, pull requests, and CI **workflow runs** (500+ rows total).

| Table             | Rows | Purpose |
|-------------------|------|---------|
| `github_users`    | 25   | Logins, display names, companies |
| `repositories`    | 20   | `full_name`, stars, visibility, branch |
| `pull_requests`   | 300  | Per-repo PR numbers, state, merge stats, authors |
| `workflow_runs`   | 220  | Linked to PRs; workflow name, event, conclusion, duration |

**Total: 565 rows** (good for `JOIN`, `GROUP BY`, filters, and aggregates).

## Load

### Docker Compose (automatic on first start)

From `tutorials/postgreSQL`:

```bash
cd tutorials/postgreSQL
docker compose up -d
```

The compose file mounts `practice/github_sample` into `/docker-entrypoint-initdb.d`, so **`github_practice_seed.sql` runs automatically** the first time Postgres creates the data directory (`devdb`).

- Connect: `psql postgresql://postgres:postgres@localhost:5432/devdb`
- Re-seed from scratch: `docker compose down -v` (wipes data), then `docker compose up -d` again.

### Manual `psql` (any Postgres)

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

## Example queries

**PRs with repo and owner**

```sql
SELECT r.full_name, pr.number, pr.title, pr.state, pr.merged, u.login AS author
FROM pull_requests pr
JOIN repositories r ON r.id = pr.repo_id
JOIN github_users u ON u.login = pr.author_login
WHERE r.full_name = 'ranjith-dev/devops'
ORDER BY pr.number;
```

**Workflow success rate by repository**

```sql
SELECT r.full_name,
       count(*) FILTER (WHERE w.conclusion = 'success') AS successes,
       count(*) AS runs,
       round(100.0 * count(*) FILTER (WHERE w.conclusion = 'success') / nullif(count(*), 0), 1) AS success_pct
FROM workflow_runs w
JOIN repositories r ON r.id = w.repo_id
GROUP BY r.id, r.full_name
ORDER BY success_pct DESC;
```

**Authors with most merged PRs**

```sql
SELECT author_login, count(*) AS merged_prs
FROM pull_requests
WHERE merged = TRUE
GROUP BY author_login
ORDER BY merged_prs DESC
LIMIT 10;
```

**Failed workflow runs with PR title**

```sql
SELECT w.workflow_name, w.conclusion, w.duration_seconds, pr.title, r.full_name
FROM workflow_runs w
JOIN pull_requests pr ON pr.id = w.pr_id
JOIN repositories r ON r.id = w.repo_id
WHERE w.conclusion = 'failure'
ORDER BY w.created_at DESC
LIMIT 20;
```

**Open PRs older than a synthetic cutoff (window practice)**

```sql
SELECT pr.id, r.full_name, pr.title, pr.created_at
FROM pull_requests pr
JOIN repositories r ON r.id = pr.repo_id
WHERE pr.state = 'open'
  AND pr.created_at < timestamptz '2024-06-01 00:00:00+00';
```

Re-run the seed script anytime to **drop and recreate** all four tables (destructive for this dataset only).

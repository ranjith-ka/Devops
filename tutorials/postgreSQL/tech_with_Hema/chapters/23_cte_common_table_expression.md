## 23. CTE (Common Table Expression)

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=11525s) · [⬆ Back to top](#table-of-contents)

A **CTE** (Common Table Expression) defines a temporary named result set with `WITH` — a named subquery you can reference in the main query. It improves readability and breaks complex logic into steps.

### Syntax

- **Single CTE:** `WITH cte_name AS (SELECT ...) SELECT ... FROM cte_name;`
- **Multiple CTEs:** `WITH cte1 AS (...), cte2 AS (...) SELECT ...;`

### Load the GitHub practice dataset

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

Examples below use `pull_requests`, `repositories`, `workflow_runs`, and (in one example) the view `pr_with_repo_and_author` from [chapter 18](18_views.md).

---

### 1) Simple CTE — filter first, then aggregate

Merged PRs only, then aggregate by author.

```sql
WITH merged_prs AS (
  SELECT
    id,
    repo_id,
    author_login,
    additions,
    deletions
  FROM pull_requests
  WHERE merged = TRUE
)
SELECT
  author_login,
  COUNT(*) AS merged_pr_count,
  COALESCE(SUM(additions), 0) AS total_additions,
  COALESCE(SUM(deletions), 0) AS total_deletions
FROM merged_prs
GROUP BY author_login
ORDER BY merged_pr_count DESC, total_additions DESC;
```

---

### 2) Multiple CTEs — step-by-step analytics

a) successful workflow runs → b) join to PRs → c) summarize by repo.

```sql
WITH successful_runs AS (
  SELECT
    repo_id,
    pr_id,
    duration_seconds
  FROM workflow_runs
  WHERE conclusion = 'success'
),
runs_with_pr AS (
  SELECT
    sr.repo_id,
    pr.id AS pr_id,
    pr.author_login,
    sr.duration_seconds
  FROM successful_runs sr
  JOIN pull_requests pr
    ON pr.id = sr.pr_id
),
per_repo AS (
  SELECT
    repo_id,
    COUNT(*) AS successful_run_count,
    AVG(duration_seconds) AS avg_duration_seconds
  FROM runs_with_pr
  GROUP BY repo_id
)
SELECT
  r.full_name,
  p.successful_run_count,
  p.avg_duration_seconds
FROM per_repo p
JOIN repositories r
  ON r.id = p.repo_id
ORDER BY p.successful_run_count DESC, p.avg_duration_seconds;
```

---

### 3) CTE for “latest per group” — `DISTINCT ON`

Latest workflow run per `(repo_id, pr_id)`.

```sql
WITH latest_run AS (
  SELECT DISTINCT ON (repo_id, pr_id)
    id,
    repo_id,
    pr_id,
    workflow_name,
    status,
    conclusion,
    created_at,
    duration_seconds
  FROM workflow_runs
  ORDER BY repo_id, pr_id, created_at DESC, id DESC
)
SELECT
  lr.repo_id,
  lr.pr_id,
  lr.workflow_name,
  lr.status,
  lr.conclusion,
  lr.created_at,
  lr.duration_seconds
FROM latest_run lr
ORDER BY lr.created_at DESC;
```

---

### 4) CTE to simplify a join-heavy query (uses the chapter 18 view)

Requires `CREATE VIEW pr_with_repo_and_author` as in [chapter 18](18_views.md). Merged PRs, then count by author company.

```sql
WITH merged_view AS (
  SELECT
    id,
    repo,
    author_company
  FROM pr_with_repo_and_author
  WHERE merged = TRUE
)
SELECT
  author_company,
  COUNT(*) AS merged_pr_count
FROM merged_view
GROUP BY author_company
ORDER BY merged_pr_count DESC, author_company;
```

---

### 5) Optional — generic pattern (orders / customers)

Illustrates a CTE as a reusable filtered set when you have `orders` and `customers` tables (not in the GitHub seed).

```sql
WITH completed_orders AS (
  SELECT
    id,
    customer_id,
    product_name,
    order_date
  FROM orders
  WHERE status = 'completed'
)
SELECT
  c.id,
  c."name",
  COUNT(co.id) AS completed_order_count
FROM customers c
LEFT JOIN completed_orders co
  ON co.customer_id = c.id
GROUP BY c.id, c."name"
ORDER BY completed_order_count DESC, c.id;
```

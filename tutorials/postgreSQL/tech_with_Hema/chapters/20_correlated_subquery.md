## 20. Correlated Subquery

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=10256s) · [⬆ Back to top](#table-of-contents)

A subquery that references columns from the outer query — executes once per outer row.

```sql
SELECT e.name, e.salary
FROM employees e
WHERE e.salary > (
    SELECT AVG(e2.salary) FROM employees e2 WHERE e2.dept_id = e.dept_id
);
```

**Contrast with a non-correlated subquery:** `WHERE salary > (SELECT AVG(salary) FROM employees)` uses one global average. The query above uses **`e.dept_id`** inside the subquery, so the average is **per department** — that reference to the outer row is what makes it **correlated**.

### Examples on the GitHub practice dataset

Load data once (PostgreSQL), e.g. `devdb`:

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

**Per-repo comparison** — PRs whose `additions` are **above the average additions for that same repo** (the subquery uses outer `pr.repo_id`):

```sql
SELECT pr.id, pr.repo_id, pr.number, pr.additions
FROM pull_requests pr
WHERE pr.additions > (
    SELECT AVG(pr2.additions)
    FROM pull_requests pr2
    WHERE pr2.repo_id = pr.repo_id
)
ORDER BY pr.repo_id, pr.number
LIMIT 20;
```

**Per-author comparison** — PRs with **more deletions than the average** for PRs by the **same author** (`pr.author_login` in the subquery):

```sql
SELECT pr.id, pr.author_login, pr.deletions
FROM pull_requests pr
WHERE pr.deletions > (
    SELECT AVG(pr2.deletions)
    FROM pull_requests pr2
    WHERE pr2.author_login = pr.author_login
)
ORDER BY pr.author_login, pr.id
LIMIT 20;
```

**Correlated `EXISTS`** — repositories that have **at least one** failed workflow run (no columns needed from the subquery except the link `w.repo_id = r.id`):

```sql
SELECT r.id, r.full_name, r.stars
FROM repositories r
WHERE EXISTS (
    SELECT 1
    FROM workflow_runs w
    WHERE w.repo_id = r.id
      AND w.conclusion = 'failure'
)
ORDER BY r.full_name;
```

The engine still **optimizes** these (often similar to joins), but logically the inner condition depends on **each outer row** — that is the correlated pattern.

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

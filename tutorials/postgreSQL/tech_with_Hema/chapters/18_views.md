## 18. Views

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=7860s) · [⬆ Back to top](#table-of-contents)

A **view** is a saved SELECT query that acts like a virtual table.

```sql
CREATE VIEW active_employees AS
SELECT * FROM employees WHERE status = 'active';

SELECT * FROM active_employees;
```

### Examples on the GitHub practice dataset

Load data once (PostgreSQL), e.g. `devdb`:

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

**Virtual “table” for a common join** — PRs with repository and author details (reuse instead of repeating the join):

```sql
DROP VIEW IF EXISTS pr_with_repo_and_author;
CREATE VIEW pr_with_repo_and_author AS
SELECT
    pr.id,
    pr.number,
    pr.title,
    pr.state,
    pr.merged,
    r.full_name AS repo,
    u.login AS author_login,
    u.display_name AS author_name,
    u.company AS author_company
FROM pull_requests pr
INNER JOIN repositories r ON r.id = pr.repo_id
INNER JOIN github_users u ON u.login = pr.author_login;

SELECT * FROM pr_with_repo_and_author LIMIT 10;
```

**Filtered slice** — only open PRs, same as querying a table:

```sql
DROP VIEW IF EXISTS open_pull_requests;
CREATE VIEW open_pull_requests AS
SELECT pr.id, pr.number, pr.title, r.full_name AS repo, pr.author_login
FROM pull_requests pr
INNER JOIN repositories r ON r.id = pr.repo_id
WHERE pr.state = 'open';

SELECT * FROM open_pull_requests ORDER BY repo, number LIMIT 15;
```

**Workflow runs that failed** (handy for dashboards or alerts):

```sql
DROP VIEW IF EXISTS failed_workflow_runs;
CREATE VIEW failed_workflow_runs AS
SELECT
    wr.id AS run_id,
    wr.workflow_name,
    wr.repo_id,
    r.full_name AS repo,
    wr.pr_id,
    pr.title AS pr_title,
    wr.duration_seconds,
    wr.created_at
FROM workflow_runs wr
INNER JOIN repositories r ON r.id = wr.repo_id
LEFT JOIN pull_requests pr ON pr.id = wr.pr_id
WHERE wr.conclusion = 'failure';

SELECT * FROM failed_workflow_runs ORDER BY created_at DESC LIMIT 20;
```

`DROP VIEW IF EXISTS` keeps re-runs of these scripts friendly while you practice.

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

## 16. Joins (Inner / Left / Right)

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=6504s) · [⬆ Back to top](#table-of-contents)

Combine rows from two or more tables based on a related column.

**Summary:** **INNER** — only rows that match on both sides. **LEFT** — every row from the first table; missing matches show NULL on the second. **RIGHT** — same as LEFT with tables swapped (rare); prefer LEFT for readability.

```sql
-- Inner join: only matching rows
SELECT e.name, d.department_name
FROM employees e
INNER JOIN departments d ON e.dept_id = d.id;

-- Left join: all from left + matching from right
SELECT e.name, d.department_name
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.id;
```

### Examples on the GitHub practice dataset

Load data once (PostgreSQL), e.g. `devdb`:

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

Tables: `pull_requests`, `repositories`, `github_users`, `workflow_runs`.

**INNER JOIN** — only rows where the join condition matches.

```sql
SELECT pr.id, pr.number, pr.title, r.full_name AS repo
FROM pull_requests pr
INNER JOIN repositories r ON r.id = pr.repo_id
ORDER BY r.full_name, pr.number
LIMIT 10;
```

```sql
SELECT pr.id, pr.state, u.login, u.display_name, u.company
FROM pull_requests pr
INNER JOIN github_users u ON u.login = pr.author_login
WHERE pr.state = 'open'
LIMIT 10;
```

**LEFT JOIN** — every row from the **left** table; right side is NULL when there is no match.

```sql
SELECT u.login, u.display_name, pr.id AS pr_id, pr.title
FROM github_users u
LEFT JOIN pull_requests pr ON pr.author_login = u.login
ORDER BY u.login
LIMIT 25;
```

```sql
SELECT r.full_name, pr.number, pr.state
FROM repositories r
LEFT JOIN pull_requests pr ON pr.repo_id = r.id
ORDER BY r.full_name, pr.number
LIMIT 15;
```

**RIGHT JOIN** — same idea as LEFT with tables swapped (less common).

```sql
SELECT r.full_name, pr.title
FROM pull_requests pr
RIGHT JOIN repositories r ON r.id = pr.repo_id
LIMIT 10;
```

**Two INNER JOINs** (PR + repo + author):

```sql
SELECT pr.id, pr.title, r.full_name, u.company
FROM pull_requests pr
INNER JOIN repositories r ON r.id = pr.repo_id
INNER JOIN github_users u ON u.login = pr.author_login
LIMIT 10;
```

**PR + workflow runs**

```sql
SELECT pr.id, pr.title, w.workflow_name, w.conclusion, w.duration_seconds
FROM pull_requests pr
INNER JOIN workflow_runs w ON w.pr_id = pr.id
LIMIT 15;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

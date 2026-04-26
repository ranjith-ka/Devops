## 15. GROUP BY / HAVING / Aggregates / ORDER BY / DISTINCT

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=5815s) · [⬆ Back to top](#table-of-contents)

```sql
-- Count employees per department
SELECT department, COUNT(*) AS total
FROM employees
GROUP BY department
HAVING COUNT(*) > 5
ORDER BY total DESC;

-- Unique departments
SELECT DISTINCT department FROM employees;
```

Aggregate functions: `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`.

---

---

## Practice

> Write your own queries below.

```sql
SELECT COUNT(*) FROM pull_requests;
SELECT SUM(additions) FROM pull_requests;
SELECT AVG(additions) FROM pull_requests;
SELECT MIN(additions) FROM pull_requests;
SELECT MAX(additions) FROM pull_requests;

SELECT MIN(created_at) AS first_pr, MAX(created_at) AS last_pr FROM pull_requests;
SELECT MIN(created_at) AS first_pr, MAX(created_at) AS last_pr FROM pull_requests GROUP BY repo_id;

-- 3) Average number of changed files per author
SELECT author_login,
       AVG(changed_files) AS avg_changed_files
FROM pull_requests
WHERE changed_files IS NOT NULL
GROUP BY author_login
ORDER BY author_login;

-- 4) Average number of changed files per author and state
SELECT author_login,
       state,
       AVG(changed_files) AS avg_changed_files
FROM pull_requests
WHERE changed_files IS NOT NULL
GROUP BY author_login, state
HAVING state = 'closed'
ORDER BY author_login;

-- 5) Unique authors
SELECT DISTINCT author_login FROM pull_requests pr ORDER BY author_login;
```

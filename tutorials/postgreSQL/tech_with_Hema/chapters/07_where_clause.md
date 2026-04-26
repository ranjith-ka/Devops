## 7. WHERE Clause

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=1858s) · [⬆ Back to top](#table-of-contents)

Filter rows using conditions:

```sql
SELECT * FROM employees WHERE age > 25;
SELECT * FROM employees WHERE department = 'IT' AND salary > 50000;
```
Operators: `=`, `!=`, `>`, `<`, `>=`, `<=`, `AND`, `OR`, `NOT`, `BETWEEN`, `IN`.


**Comparison operators** — compare a column (or expression) to a value:

| Operator | Meaning | Example |
|----------|---------|---------|
| `=` | equal | `WHERE department = 'IT'` |
| `!=` or `<>` | not equal | `WHERE name != 'Ravi'` |
| `>`, `<`, `>=`, `<=` | ordering | `WHERE age > 25`, `WHERE salary <= 50000` |

**Logical operators** — combine conditions; each part must be a valid condition (e.g. `column op value`), not a bare string:

- **`AND`** — all conditions must be true: `WHERE age > 25 AND department = 'IT'`.
- **`OR`** — at least one condition is true: `WHERE age < 20 OR age > 60`.
- **`NOT`** — negates the next condition: `WHERE NOT (age BETWEEN 30 AND 50)`.

**Shortcut operators:**

- **`BETWEEN`** — inclusive range: `WHERE age BETWEEN 18 AND 65` (same as `age >= 18 AND age <= 65`).
- **`IN`** — value in a list: `WHERE department IN ('IT', 'HR', 'Finance')`.

Parentheses control order when mixing `AND` / `OR`: `WHERE (dept = 'IT' OR dept = 'HR') AND salary > 40000`.

---

---

## Practice

Uses the GitHub practice dataset (`github_users`, `repositories`, `pull_requests`). Load once if needed:

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

> Write your own queries below.

```sql
-- IN — logins (same idea as name IN (...))
SELECT * FROM github_users WHERE login IN ('amit-v', 'ravi-p', 'priya-m');

-- BETWEEN — inclusive range on numeric column (here: lines added in a PR)
SELECT * FROM pull_requests WHERE additions BETWEEN 100 AND 300;

-- AND — both conditions
SELECT * FROM repositories WHERE owner_login = 'ranjith-dev' AND stars > 100;

-- AND — combine comparison + categorical filter
SELECT * FROM pull_requests WHERE additions > 200 AND state = 'closed';

-- OR — either condition
SELECT * FROM pull_requests WHERE additions < 80 OR additions > 750;

-- NOT — negate a BETWEEN
SELECT * FROM pull_requests WHERE NOT (additions BETWEEN 150 AND 500);

-- IN — match any of several values
SELECT * FROM github_users WHERE company IN ('Acme Cloud', 'StartupXYZ', 'FinTech Co');

-- Parentheses — (A OR B) AND C
SELECT * FROM repositories
WHERE (owner_login = 'ranjith-dev' OR owner_login = 'hema-s') AND stars > 50;
```

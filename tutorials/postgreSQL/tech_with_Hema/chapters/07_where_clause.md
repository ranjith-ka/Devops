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

> Write your own queries below.

```sql
-- Your practice queries here

SELECT * FROM employees WHERE name IN ('Amit', 'Ravi', 'Priya');
SELECT * FROM employees WHERE age BETWEEN 20 AND 30;
SELECT * FROM employees WHERE department = 'IT' AND salary > 50000;
SELECT * FROM employees WHERE age > 25 AND department = 'IT';
SELECT * FROM employees WHERE age < 20 OR age > 60;
SELECT * FROM employees WHERE NOT (age BETWEEN 30 AND 50);
SELECT * FROM employees WHERE department IN ('IT', 'HR', 'Finance');
SELECT * FROM employees WHERE (dept = 'IT' OR dept = 'HR') AND salary > 40000;
```

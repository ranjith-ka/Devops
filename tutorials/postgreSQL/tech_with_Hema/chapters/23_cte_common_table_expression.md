## 23. CTE (Common Table Expression)

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=11525s) · [⬆ Back to top](#table-of-contents)

A temporary named result set for cleaner, modular queries.

```sql
WITH high_earners AS (
    SELECT * FROM employees WHERE salary > 80000
)
SELECT department, COUNT(*) FROM high_earners GROUP BY department;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

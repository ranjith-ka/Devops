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

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

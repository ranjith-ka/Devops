## 24. Aggregate Window Functions

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=12408s) · [⬆ Back to top](#table-of-contents)

Perform calculations across a set of rows **without collapsing** them.

```sql
SELECT name, department, salary,
       SUM(salary) OVER (PARTITION BY department) AS dept_total
FROM employees;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

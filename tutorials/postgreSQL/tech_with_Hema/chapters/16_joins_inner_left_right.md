## 16. Joins (Inner / Left / Right)

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=6504s) · [⬆ Back to top](#table-of-contents)

Combine rows from two or more tables based on a related column.

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

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

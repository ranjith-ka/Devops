## 17. Self Join / Cross Join

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=7246s) · [⬆ Back to top](#table-of-contents)

- **Self join**: a table joined with itself (e.g. employee → manager).
- **Cross join**: every row from table A combined with every row from table B (cartesian product).

```sql
-- Self join: find each employee's manager
SELECT e.name AS employee, m.name AS manager
FROM employees e
JOIN employees m ON e.manager_id = m.id;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

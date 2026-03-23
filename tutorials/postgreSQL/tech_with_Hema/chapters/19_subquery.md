## 19. Subquery

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=9031s) · [⬆ Back to top](#table-of-contents)

A query nested inside another query.

```sql
SELECT name FROM employees
WHERE salary > (SELECT AVG(salary) FROM employees);
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

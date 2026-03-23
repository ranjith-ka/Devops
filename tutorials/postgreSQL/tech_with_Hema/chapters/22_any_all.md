## 22. ANY / ALL

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=10912s) · [⬆ Back to top](#table-of-contents)

Compare a value against a set returned by a subquery.

```sql
-- Salary greater than ANY in IT department
SELECT name FROM employees
WHERE salary > ANY (SELECT salary FROM employees WHERE department = 'IT');
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

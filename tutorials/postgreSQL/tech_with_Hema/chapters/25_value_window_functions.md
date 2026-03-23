## 25. Value Window Functions

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=13025s) · [⬆ Back to top](#table-of-contents)

Access values from other rows: `FIRST_VALUE`, `LAST_VALUE`, `LEAD`, `LAG`.

```sql
SELECT name, salary,
       LAG(salary) OVER (ORDER BY salary) AS prev_salary,
       LEAD(salary) OVER (ORDER BY salary) AS next_salary
FROM employees;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

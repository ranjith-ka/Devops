## 15. GROUP BY / HAVING / Aggregates / ORDER BY / DISTINCT

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=5815s) · [⬆ Back to top](#table-of-contents)

```sql
-- Count employees per department
SELECT department, COUNT(*) AS total
FROM employees
GROUP BY department
HAVING COUNT(*) > 5
ORDER BY total DESC;

-- Unique departments
SELECT DISTINCT department FROM employees;
```

Aggregate functions: `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`.

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

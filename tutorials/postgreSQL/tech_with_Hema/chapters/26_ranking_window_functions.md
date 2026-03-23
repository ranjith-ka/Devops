## 26. Ranking Window Functions

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=13593s) · [⬆ Back to top](#table-of-contents)

Rank rows: `ROW_NUMBER`, `RANK`, `DENSE_RANK`, `NTILE`.

```sql
SELECT name, salary,
       ROW_NUMBER() OVER (ORDER BY salary DESC) AS row_num,
       RANK()       OVER (ORDER BY salary DESC) AS rnk,
       DENSE_RANK() OVER (ORDER BY salary DESC) AS dense_rnk
FROM employees;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

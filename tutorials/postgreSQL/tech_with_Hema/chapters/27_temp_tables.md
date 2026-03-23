## 27. Temp Tables

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=14322s) · [⬆ Back to top](#table-of-contents)

Temporary storage for intermediate results within a session.

```sql
CREATE TABLE #temp_employees (id INT, name VARCHAR(100));
INSERT INTO #temp_employees SELECT id, name FROM employees WHERE dept_id = 1;
SELECT * FROM #temp_employees;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

## 28. Stored Procedures

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=15048s) · [⬆ Back to top](#table-of-contents)

Reusable blocks of SQL with parameters and control flow.

```sql
CREATE PROCEDURE GetEmployeesByDept @dept_id INT
AS
BEGIN
    SELECT * FROM employees WHERE dept_id = @dept_id;
END;

EXEC GetEmployeesByDept @dept_id = 3;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

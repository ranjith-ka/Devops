## 31. Functions (Scalar / Table-Valued)

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=18169s) · [⬆ Back to top](#table-of-contents)

- **Scalar function**: returns a single value.
- **Table-valued function**: returns a table.

```sql
-- Scalar
CREATE FUNCTION dbo.GetFullName(@first VARCHAR(50), @last VARCHAR(50))
RETURNS VARCHAR(101)
AS BEGIN RETURN @first + ' ' + @last; END;

-- Table-valued
CREATE FUNCTION dbo.GetDeptEmployees(@dept_id INT)
RETURNS TABLE
AS RETURN (SELECT * FROM employees WHERE dept_id = @dept_id);
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

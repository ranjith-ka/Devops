## 39. Dynamic SQL

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=27184s) · [⬆ Back to top](#table-of-contents)

Build and execute SQL strings at runtime.

```sql
DECLARE @table NVARCHAR(100) = N'employees';
DECLARE @sql NVARCHAR(MAX) = N'SELECT * FROM ' + QUOTENAME(@table);
EXEC sp_executesql @sql;
```

> Always use `QUOTENAME` or parameterized queries with `sp_executesql` to prevent SQL injection.

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

## 32. Cursors

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=19664s) · [⬆ Back to top](#table-of-contents)

Process rows one at a time when set-based operations are insufficient.

```sql
DECLARE @name VARCHAR(100);
DECLARE emp_cursor CURSOR FOR SELECT name FROM employees;
OPEN emp_cursor;
FETCH NEXT FROM emp_cursor INTO @name;
WHILE @@FETCH_STATUS = 0
BEGIN
    PRINT @name;
    FETCH NEXT FROM emp_cursor INTO @name;
END;
CLOSE emp_cursor;
DEALLOCATE emp_cursor;
```

> **Note**: Prefer set-based operations over cursors for performance.

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

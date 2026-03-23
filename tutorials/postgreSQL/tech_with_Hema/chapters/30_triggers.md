## 30. Triggers

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=17053s) · [⬆ Back to top](#table-of-contents)

Automatic actions that fire on INSERT, UPDATE, or DELETE.

```sql
CREATE TRIGGER trg_audit
ON employees
AFTER UPDATE
AS
BEGIN
    INSERT INTO audit_log (action, timestamp)
    VALUES ('UPDATE', GETDATE());
END;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

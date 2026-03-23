## 33. Transactions

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=21250s) · [⬆ Back to top](#table-of-contents)

Group operations so they all succeed or all fail together.

```sql
BEGIN TRANSACTION;
    UPDATE accounts SET balance = balance - 500 WHERE id = 1;
    UPDATE accounts SET balance = balance + 500 WHERE id = 2;
COMMIT;
-- Use ROLLBACK instead of COMMIT to undo on error.
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

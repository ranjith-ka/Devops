## 21. EXISTS / NOT EXISTS

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=10675s) · [⬆ Back to top](#table-of-contents)

Check whether a subquery returns any rows.

```sql
SELECT name FROM customers c
WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id);
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

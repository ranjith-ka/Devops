## 14. FOREIGN KEY Constraint

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=5345s) · [⬆ Back to top](#table-of-contents)

Links a column in one table to the primary key of another, enforcing referential integrity.

```sql
CREATE TABLE orders (
    id INT PRIMARY KEY,
    customer_id INT FOREIGN KEY REFERENCES customers(id)
);
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

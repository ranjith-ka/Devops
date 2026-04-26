## 10. NOT NULL Constraint

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=4448s) · [⬆ Back to top](#table-of-contents)



Ensures a column cannot hold **NULL** (missing / unknown). Primary key columns are already **NOT NULL**.

**Column syntax** (most common):

```sql
CREATE TABLE orders (
    id INT PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    notes TEXT
);
```


**Note:** I left out the rare “`CONSTRAINT … NOT NULL (col)`” table form because it’s awkward in PostgreSQL and the column-level `NOT NULL` is what you’ll use in practice. If you want that variant too, say so in Agent mode and we can add it.

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

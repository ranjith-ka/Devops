## 14. FOREIGN KEY Constraint

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=5345s) · [⬆ Back to top](#table-of-contents)

Links a column in one table to the primary key of another, enforcing referential integrity.

In simple words:

- The child table stores an id that points to a row in the parent table.
- You can only use values that **already exist** in the parent (for that column).
- This stops “broken links”: for example, an order cannot point to a customer that is not in the database.
- The database checks this for you when you **insert**, **update**, or **delete** rows.

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
-- 1) Parent table (must exist before the FK)
CREATE TABLE customers (
    id          INT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(255) UNIQUE
);

-- 2) Child table: FK using column-level REFERENCES
CREATE TABLE orders (
    id              INT PRIMARY KEY,
    customer_id     INT REFERENCES customers (id),
    product_name    VARCHAR(100) NOT NULL,
    order_date      DATE DEFAULT CURRENT_DATE,
    status          VARCHAR(20) DEFAULT 'pending'
);

-- Example data (parent first, then child)
INSERT INTO customers (id, name, email) VALUES
    (1, 'Acme Corp', 'contact@acme.test');

INSERT INTO orders (id, customer_id, product_name) VALUES
    (100, 1, 'Widget');


```error
devdb=# INSERT INTO orders (id, customer_id, product_name) VALUES
    (100, 2, 'Widget');
ERROR:  foreign key constraint "orders_customer_id_fkey" on table "orders" may not be NULL
LINE 1: INSERT INTO orders (id, customer_id, product_name) VALUES
                    ^
devdb=# 
```

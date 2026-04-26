## 13. DEFAULT Constraint

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=5189s) · [⬆ Back to top](#table-of-contents)

Assigns a default value when no value is provided.

```sql
CREATE TABLE tasks (
    id INT PRIMARY KEY,
    status VARCHAR(20) DEFAULT 'pending'
);
```

---

---

## Practice

> Write your own queries below.

```sql

CREATE TABLE orders (
    id INT PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    notes TEXT,
    email VARCHAR(255) UNIQUE,
    order_date DATE CHECK (order_date >= DATE '1980-01-01'),
    status VARCHAR(20) DEFAULT 'pending'
    
);


INSERT INTO orders (id, product_name, notes, email, order_date)
VALUES (15, 'Keyboard', 'Mechanical', 'dg@fg.com', '2002-01-01');
```

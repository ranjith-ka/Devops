## 6. Basic SQL Commands

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=1078s) · [⬆ Back to top](#table-of-contents)

Hands-on practice with the foundational SQL statements.

**Setup** (run in `psql` once so `employees` exists and `SELECT` returns rows before you add Ravi):

```sql
DROP TABLE IF EXISTS employees;
CREATE TABLE employees (
    name varchar(100) NOT NULL,
    age int
);

INSERT INTO employees (name, age) VALUES
    ('Priya', 32),
    ('Amit', 41);
```

**Example sequence:**

```sql
SELECT * FROM employees;
INSERT INTO employees (name, age) VALUES ('Ravi', 28);
UPDATE employees SET age = 29 WHERE name = 'Ravi';
DELETE FROM employees WHERE name = 'Ravi';
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

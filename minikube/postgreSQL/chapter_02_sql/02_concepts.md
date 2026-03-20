# 2.2 Concepts

Official: [2.2. Concepts](https://www.postgresql.org/docs/current/tutorial-concepts.html)

PostgreSQL is a **relational** DBMS: data lives in **relations** (tables).

- **Table:** named collection of **rows**.
- **Row:** same set of **named columns**; each column has a **data type**.
- Column order is fixed per row; **row order** in a table is **not** guaranteed by SQL unless you use `ORDER BY`.

**Database cluster:** one PostgreSQL server instance manages many **databases**; tables belong to a database.

Next: [03_create_table.md](03_create_table.md)

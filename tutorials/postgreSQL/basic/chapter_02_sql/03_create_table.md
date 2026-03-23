# 2.3 Creating a New Table

Official: [2.3. Creating a New Table](https://www.postgresql.org/docs/current/tutorial-table.html)

Define **name**, **columns**, and **types**. Statements end with **`;`**. Whitespace is flexible. **`--`** starts a line comment.

SQL keywords and unquoted identifiers are **case-insensitive**; double-quoted identifiers preserve case.

```sql
CREATE TABLE weather (
    city            varchar(80),
    temp_lo         int,           -- low temperature
    temp_hi         int,           -- high temperature
    prcp            real,          -- precipitation
    date            date
);
```

Types used here: `varchar(80)`, `int`, `real`, `date`. PostgreSQL also supports `smallint`, `double precision`, `time`, `timestamp`, `interval`, geometric types, and **user-defined** types.

```sql
CREATE TABLE cities (
    name            varchar(80),
    location        point
);
```

`point` is PostgreSQL-specific.

Remove a table:

```sql
DROP TABLE tablename;
```

Prev: [02_concepts.md](02_concepts.md) · Next: [04_populate.md](04_populate.md)

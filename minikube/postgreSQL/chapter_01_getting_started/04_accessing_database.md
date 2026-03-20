# 1.4 Accessing a database

Official: [1.4. Accessing a Database](https://www.postgresql.org/docs/current/tutorial-accessdb.html)

Ways to use a database:

1. **`psql`** — interactive terminal (what the tutorial uses).  
2. **GUI tools** — e.g. pgAdmin; ODBC/JDBC in spreadsheets or BI tools (not covered in the core tutorial).  
3. **Applications** — language drivers and APIs ([Part IV. Client Interfaces](https://www.postgresql.org/docs/current/client-interfaces.html)).

## psql

```bash
psql mydb
```

Omit the database name and `psql` connects to a database named like your **OS user** (same idea as bare `createdb`).

### Prompt

- **`mydb=>`** — normal role  
- **`mydb=#`** — **superuser** (common if you installed the instance yourself); not subject to normal RLS-style access rules in the same way

### Try SQL

```sql
SELECT version();
SELECT current_date;
SELECT 2 + 2;
```

### psql backslash commands (not SQL)

| Command | Purpose |
|---------|--------|
| `\h` | Help for SQL commands |
| `\?` | Help for psql meta-commands |
| `\q` | Quit |

Full reference: [psql](https://www.postgresql.org/docs/current/app-psql.html).

### If psql fails

Diagnostics overlap with **`createdb`**. If `createdb` worked against the same host/port/user, `psql` should too.

## Docker (this repo)

```bash
psql -h localhost -p 5432 -U postgres -d devdb
```

Or after creating **`mydb`**:

```bash
psql -h localhost -p 5432 -U postgres -d mydb
```

Prev: [03_creating_database.md](03_creating_database.md) · Next: [Chapter 2 — SQL](../chapter_02_sql/README.md)

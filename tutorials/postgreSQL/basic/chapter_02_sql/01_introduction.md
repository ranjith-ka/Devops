# 2.1 Introduction

Official: [2.1. Introduction](https://www.postgresql.org/docs/current/tutorial-sql-intro.html)

This chapter shows **simple SQL** in PostgreSQL. It is an introduction only, not a full SQL course. PostgreSQL includes **standard SQL** plus **extensions**.

**Assumption:** You have a database `mydb` and can start `psql` (see Chapter 1).

## Optional: tutorial files from PostgreSQL source

Examples also live under `src/tutorial/` in the PostgreSQL source. There:

```bash
cd .../src/tutorial
make
psql -s mydb
```

In `psql`, `\i basics.sql` loads the script. The `-s` flag enables single-step mode (pause before each statement).

- **`\i`** — read commands from a file.

Next: [02_concepts.md](02_concepts.md)

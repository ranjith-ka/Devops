# 1.3 Creating a database

Official: [1.3. Creating a Database](https://www.postgresql.org/docs/current/tutorial-createdb.html)

One PostgreSQL server holds **many databases** (often one per project or user). If an admin already created one for you, skip to [04_accessing_database.md](04_accessing_database.md).

## Create and drop

```bash
createdb mydb
```

No output usually means success.

**Default name:** `createdb` with no argument creates a database named like your **OS user** (many tools default to that).

```bash
dropdb mydb
```

You **must** pass the name to `dropdb`. Dropping removes files on the server; it is **not reversible**.

## Common errors

| Message | Meaning |
|---------|--------|
| `createdb: command not found` | PostgreSQL client tools not on `PATH` (try full path, e.g. `/usr/local/pgsql/bin/createdb`). |
| `connection ... failed ... Is the server running` | Server down or wrong host/port/socket. |
| `FATAL: role "you" does not exist` | No matching **PostgreSQL role** (login). Roles are **not** the same as OS users. Fix with admin, or `-U` / `PGUSER` for the right role. |
| `permission denied to create database` | Your role lacks privilege; admin must grant `CREATEDB` or you connect as a user that can create DBs. |

**Note:** The OS user that **started** the server typically has a same-named superuser role that **can** create databases (see official footnote on [createdb](https://www.postgresql.org/docs/current/tutorial-createdb.html)).

## Docker (this repo)

Server may already have **`devdb`**. For tutorial name **`mydb`**:

```bash
createdb -h localhost -p 5432 -U postgres mydb
```

Password: set via `PGPASSWORD` or a `.pgpass` file, or enter when prompted.

Prev: [02_architecture.md](02_architecture.md) · Next: [04_accessing_database.md](04_accessing_database.md)

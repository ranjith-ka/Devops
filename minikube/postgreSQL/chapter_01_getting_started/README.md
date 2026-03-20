# Chapter 1 — Getting Started

Official: [Chapter 1. Getting Started](https://www.postgresql.org/docs/current/tutorial-start.html) (PostgreSQL docs).

## Sections (this folder)

| # | Topic | Local | Official |
|---|--------|-------|----------|
| 1 | Installation | [01_installation.md](01_installation.md) | [1.1](https://www.postgresql.org/docs/current/tutorial-install.html) |
| 2 | Architectural fundamentals | [02_architecture.md](02_architecture.md) | [1.2](https://www.postgresql.org/docs/current/tutorial-arch.html) |
| 3 | Creating a database | [03_creating_database.md](03_creating_database.md) | [1.3](https://www.postgresql.org/docs/current/tutorial-createdb.html) |
| 4 | Accessing a database | [04_accessing_database.md](04_accessing_database.md) | [1.4](https://www.postgresql.org/docs/current/tutorial-accessdb.html) |

## Next in this repo

[Chapter 2 — SQL language](../chapter_02_sql/README.md) (official [Chapter 2](https://www.postgresql.org/docs/current/tutorial-sql.html)).

## PostgreSQL in Docker (this repo)

From repo root:

```bash
docker compose -f minikube/postgreSQL/docker-compose.yaml up -d
```

Connect (default DB **`devdb`**, user/password **`postgres`**):

```bash
psql -h localhost -p 5432 -U postgres -d devdb
```

The tutorial often uses **`mydb`**. Create it once:

```bash
createdb -h localhost -p 5432 -U postgres mydb
# or inside psql: CREATE DATABASE mydb;
```

Then: `psql -h localhost -p 5432 -U postgres -d mydb`.

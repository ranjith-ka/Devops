# Chapter 2 — The SQL Language

Official: [Chapter 2. The SQL Language](https://www.postgresql.org/docs/current/tutorial-sql.html) (PostgreSQL docs).

**Prerequisites:** Database `mydb` (or your choice) and `psql` as in [Getting Started](https://www.postgresql.org/docs/current/tutorial-start.html).

**Hands-on:** [EXERCISES.md](EXERCISES.md) — step-by-step practice (with collapsible solutions).

## Sections (this folder)

| # | Topic | Local | Official |
|---|--------|-------|----------|
| 1 | Introduction | [01_introduction.md](01_introduction.md) | [2.1](https://www.postgresql.org/docs/current/tutorial-sql-intro.html) |
| 2 | Concepts | [02_concepts.md](02_concepts.md) | [2.2](https://www.postgresql.org/docs/current/tutorial-concepts.html) |
| 3 | Creating a table | [03_create_table.md](03_create_table.md) | [2.3](https://www.postgresql.org/docs/current/tutorial-table.html) |
| 4 | Populating rows | [04_populate.md](04_populate.md) | [2.4](https://www.postgresql.org/docs/current/tutorial-populate.html) |
| 5 | Querying | [05_query.md](05_query.md) | [2.5](https://www.postgresql.org/docs/current/tutorial-select.html) |
| 6 | Joins | [06_joins.md](06_joins.md) | [2.6](https://www.postgresql.org/docs/current/tutorial-join.html) |
| 7 | Aggregates | [07_aggregates.md](07_aggregates.md) | [2.7](https://www.postgresql.org/docs/current/tutorial-agg.html) |
| 8 | Updates | [08_updates.md](08_updates.md) | [2.8](https://www.postgresql.org/docs/current/tutorial-update.html) |
| 9 | Deletions | [09_deletions.md](09_deletions.md) | [2.9](https://www.postgresql.org/docs/current/tutorial-delete.html) |

## Run all examples from one script

From the repo:

```bash
psql -d mydb -f minikube/postgreSQL/chapter_02_sql/basics.sql
```

Or inside `psql`:

```text
\i path/to/chapter_02_sql/basics.sql
```

The upstream tutorial also ships `src/tutorial/basics.sql` in the PostgreSQL source tree; see [2.1 Introduction](https://www.postgresql.org/docs/current/tutorial-sql-intro.html).

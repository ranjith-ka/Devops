## 9. Primary Key Constraint

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=3888s) · [⬆ Back to top](#table-of-contents)

### The idea in plain language

Each row in a table needs a **stable “address”** so the database can point at exactly one row. That address is the **primary key (PK)**.

| Rule | Meaning |
|------|--------|
| **Unique** | Two rows cannot share the same PK value. |
| **Not NULL** | Every row must have a PK value. |
| **One per table** | You declare **one** primary key (one column, or several columns **together**). |

In your GitHub practice data, **`workflow_runs`** uses a simple numeric key: **`id`**. Every CI/build row gets its own **`id`**; `repo_id`, `workflow_name`, etc. describe the run but **do not replace** that identity unless you designed the table differently.

The line that defines the PK in the seed looks like this (the rest of the columns are just details):

```sql
-- From workflow_runs in github_practice_seed.sql
id SERIAL PRIMARY KEY,
```

- **`SERIAL`** — PostgreSQL fills in `1, 2, 3, …` for you when you insert (unless you supply `id` yourself).
- **`PRIMARY KEY`** — tells the database: this column is the row’s unique address.

### Two shapes (know these names)

1. **Single-column PK** — e.g. `id` alone (what **`workflow_runs`** uses).
2. **Composite PK** — several columns **together** form one key, e.g. `PRIMARY KEY (repo_id, workflow_name, run_attempt)`. Use when **no single column** is enough to say “this row” uniquely.

---

### Try it (safe scratch table)

This uses **`lesson_workflow_runs`** so you never touch the seeded **`workflow_runs`** table. Optional: load the full sample first.

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

**1) Duplicate `id` is not allowed**

```sql
DROP TABLE IF EXISTS lesson_workflow_runs CASCADE;

CREATE TABLE lesson_workflow_runs (
    id              SERIAL PRIMARY KEY,
    workflow_name   VARCHAR(200) NOT NULL
);

INSERT INTO lesson_workflow_runs (workflow_name) VALUES ('CI');
INSERT INTO lesson_workflow_runs (workflow_name) VALUES ('Lint');
-- Rows got id = 1 and id = 2 automatically.

-- This fails: id 1 is already taken
INSERT INTO lesson_workflow_runs (id, workflow_name) VALUES (1, 'Deploy');
```

You should see an error like: `duplicate key value violates unique constraint "…_pkey"`.

**2) Composite PK (same combination cannot appear twice)**

```sql
DROP TABLE IF EXISTS lesson_workflow_run_keys CASCADE;

CREATE TABLE lesson_workflow_run_keys (
    repo_id         INT NOT NULL,
    workflow_name   VARCHAR(200) NOT NULL,
    run_attempt     SMALLINT NOT NULL,
    PRIMARY KEY (repo_id, workflow_name, run_attempt)
);

INSERT INTO lesson_workflow_run_keys VALUES (1, 'CI', 1);
-- Fails: same (1, 'CI', 1) again
INSERT INTO lesson_workflow_run_keys VALUES (1, 'CI', 1);
```

---

**Takeaway:** In **`workflow_runs`**, the PK is **`id`**. Composite PKs are for when **one column is not enough** to identify a row uniquely; **`workflow_runs`** keeps life simple with **`id`** and uses foreign keys (`repo_id`, `pr_id`) to link to other tables.

> Write your own experiments below.

```sql
-- Your practice queries here

```

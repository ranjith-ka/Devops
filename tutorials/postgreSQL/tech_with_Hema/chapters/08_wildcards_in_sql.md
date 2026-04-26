## 8. Wildcards in SQL

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=2834s) · [⬆ Back to top](#table-of-contents)

Pattern matching with `LIKE` (string columns):

- **`%`** — any sequence of characters (including empty).
- **`_`** — exactly one character.

```sql
-- Starts with 'r' (e.g. ranjith-dev, ravi-p)
SELECT * FROM github_users WHERE login LIKE 'r%';

-- Ends with '-dev'
SELECT * FROM github_users WHERE login LIKE '%-dev';

-- Second character is 'a' (e.g. fatima-z, carlos-m)
SELECT * FROM github_users WHERE login LIKE '_a%';
```

More examples on the same data:

```sql
-- PR titles containing "feat"
SELECT id, title FROM pull_requests WHERE title LIKE '%feat:%';

-- Repository full name under an owner (slash is literal)
SELECT full_name FROM repositories WHERE full_name LIKE 'ranjith-dev/%';

-- Workflow names starting with "C" (e.g. CI)
SELECT DISTINCT workflow_name FROM workflow_runs WHERE workflow_name LIKE 'C%';
```

---

## Practice

Uses the GitHub practice dataset. Load once if needed:

```bash
psql -d devdb -f tutorials/postgreSQL/practice/github_sample/github_practice_seed.sql
```

> Write your own queries below.

```sql
-- Ends with "bot"
SELECT * FROM github_users WHERE login LIKE '%bot';

-- Login contains "dev" anywhere
SELECT * FROM github_users WHERE login LIKE '%dev%';

-- PR title starts with "fix:"
SELECT * FROM pull_requests WHERE title LIKE 'fix:%';

-- Branch: literal "feature/", then anything (all rows in this seed use this prefix)
SELECT * FROM pull_requests WHERE head_branch LIKE 'feature/%';

-- Workflow name "CI" — one character + "I" (underscore matches "C")
SELECT * FROM workflow_runs WHERE workflow_name LIKE '_I';
```

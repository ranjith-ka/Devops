## 38. Recursive CTE

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=25653s) · [⬆ Back to top](#table-of-contents)

Query hierarchical data (org charts, category trees) recursively.

```sql
WITH org AS (
    SELECT id, name, manager_id, 0 AS level
    FROM employees WHERE manager_id IS NULL
    UNION ALL
    SELECT e.id, e.name, e.manager_id, o.level + 1
    FROM employees e JOIN org o ON e.manager_id = o.id
)
SELECT * FROM org;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

## 34. MERGE Statement

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=22246s) · [⬆ Back to top](#table-of-contents)

Perform INSERT, UPDATE, and DELETE in a single synchronized statement.

```sql
MERGE INTO target_table AS t
USING source_table AS s ON t.id = s.id
WHEN MATCHED THEN UPDATE SET t.name = s.name
WHEN NOT MATCHED THEN INSERT (id, name) VALUES (s.id, s.name)
WHEN NOT MATCHED BY SOURCE THEN DELETE;
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

```

# 2.9 Deletions

Official: [2.9. Deletions](https://www.postgresql.org/docs/current/tutorial-delete.html)

Remove rows with **`DELETE`** and a **`WHERE`** clause:

```sql
DELETE FROM weather WHERE city = 'Hayward';
```

**Danger:** unqualified delete empties the table with **no confirmation**:

```sql
DELETE FROM tablename;
```

Always double-check `WHERE` before running deletes in production.

Prev: [08_updates.md](08_updates.md) · Next (docs): [Chapter 3. Advanced Features](https://www.postgresql.org/docs/current/tutorial-advanced.html)

# 2.8 Updates

Official: [2.8. Updates](https://www.postgresql.org/docs/current/tutorial-update.html)

Change existing rows with **`UPDATE`** + **`SET`** + optional **`WHERE`**:

```sql
UPDATE weather
    SET temp_hi = temp_hi - 2, temp_lo = temp_lo - 2
    WHERE date > '1994-11-28';
```

Verify:

```sql
SELECT * FROM weather;
```

Without `WHERE`, **every row** in the table is updated.

Prev: [07_aggregates.md](07_aggregates.md) · Next: [09_deletions.md](09_deletions.md)

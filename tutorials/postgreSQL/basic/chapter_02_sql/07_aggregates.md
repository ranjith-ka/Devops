# 2.7 Aggregate Functions

Official: [2.7. Aggregate Functions](https://www.postgresql.org/docs/current/tutorial-agg.html)

Aggregates collapse many rows into one value: `count`, `sum`, `avg`, `max`, `min`, etc.

```sql
SELECT max(temp_lo) FROM weather;
```

**You cannot use an aggregate in `WHERE`** — `WHERE` is evaluated **before** aggregation. Use a **subquery** instead:

```sql
-- WRONG: WHERE temp_lo = max(temp_lo)
SELECT city FROM weather
    WHERE temp_lo = (SELECT max(temp_lo) FROM weather);
```

## GROUP BY

One output row per group:

```sql
SELECT city, count(*), max(temp_lo)
    FROM weather
    GROUP BY city;
```

## HAVING

Filter **after** grouping (often uses aggregates):

```sql
SELECT city, count(*), max(temp_lo)
    FROM weather
    GROUP BY city
    HAVING max(temp_lo) < 40;
```

## WHERE vs HAVING

- **WHERE** — picks rows **before** groups/aggregates.
- **HAVING** — picks **groups** after aggregates.

Prefer non-aggregate filters in **WHERE** (cheaper).

```sql
SELECT city, count(*), max(temp_lo)
    FROM weather
    WHERE city LIKE 'S%'
    GROUP BY city;
```

[`LIKE`](https://www.postgresql.org/docs/current/functions-matching.html) — pattern matching (Section 9.7).

## FILTER (per aggregate)

```sql
SELECT city, count(*) FILTER (WHERE temp_lo < 45), max(temp_lo)
    FROM weather
    GROUP BY city;
```

`FILTER` limits input to **that** aggregate only; other aggregates still see all rows in the group.

Prev: [06_joins.md](06_joins.md) · Next: [08_updates.md](08_updates.md)

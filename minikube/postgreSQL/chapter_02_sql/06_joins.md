# 2.6 Joins Between Tables

Official: [2.6. Joins Between Tables](https://www.postgresql.org/docs/current/tutorial-join.html)

**Join** queries combine rows from two (or more) tables according to a condition.

## Inner join

```sql
SELECT * FROM weather JOIN cities ON city = name;
```

Unmatched rows (e.g. Hayward with no row in `cities`) **drop out**. Prefer explicit columns instead of `*`:

```sql
SELECT city, temp_lo, temp_hi, prcp, date, location
    FROM weather JOIN cities ON city = name;
```

If names collide, **qualify** with the table name:

```sql
SELECT weather.city, weather.temp_lo, weather.temp_hi,
       weather.prcp, weather.date, cities.location
    FROM weather JOIN cities ON weather.city = cities.name;
```

## Comma + WHERE (legacy)

Same idea as inner join:

```sql
SELECT *
    FROM weather, cities
    WHERE city = name;
```

Explicit `JOIN ... ON` separates join logic from other filters in `WHERE`.

## Left outer join

Keep every row from the **left** table; use NULLs when there is no right match:

```sql
SELECT *
    FROM weather LEFT OUTER JOIN cities ON weather.city = cities.name;
```

Also: `RIGHT OUTER`, `FULL OUTER` (try them).

## Self-join

Alias the same table twice:

```sql
SELECT w1.city, w1.temp_lo AS low, w1.temp_hi AS high,
       w2.city, w2.temp_lo AS low, w2.temp_hi AS high
    FROM weather w1 JOIN weather w2
        ON w1.temp_lo < w2.temp_lo AND w1.temp_hi > w2.temp_hi;
```

Aliases elsewhere:

```sql
SELECT *
    FROM weather w JOIN cities c ON w.city = c.name;
```

Prev: [05_query.md](05_query.md) · Next: [07_aggregates.md](07_aggregates.md)

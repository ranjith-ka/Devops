# 2.5 Querying a Table

Official: [2.5. Querying a Table](https://www.postgresql.org/docs/current/tutorial-select.html)

`SELECT` has: **select list** (columns/expressions), **FROM** (tables), optional **WHERE** (filter).

```sql
SELECT * FROM weather;
-- same columns explicitly:
SELECT city, temp_lo, temp_hi, prcp, date FROM weather;
```

`*` is handy ad hoc; in production, listing columns is usually preferred so schema changes do not silently change results.

## Expressions and labels

```sql
SELECT city, (temp_hi + temp_lo) / 2 AS temp_avg, date FROM weather;
```

`AS` labels the output column (optional).

## WHERE

```sql
SELECT * FROM weather
    WHERE city = 'San Francisco' AND prcp > 0.0;
```

Boolean operators: `AND`, `OR`, `NOT`.

## ORDER BY

```sql
SELECT * FROM weather ORDER BY city;
SELECT * FROM weather ORDER BY city, temp_lo;
```

## DISTINCT

```sql
SELECT DISTINCT city FROM weather;
SELECT DISTINCT city FROM weather ORDER BY city;
```

`DISTINCT` does **not** guarantee sort order unless you add `ORDER BY`.

Prev: [04_populate.md](04_populate.md) · Next: [06_joins.md](06_joins.md)

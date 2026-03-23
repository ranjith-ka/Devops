# 2.4 Populating a Table With Rows

Official: [2.4. Populating a Table With Rows](https://www.postgresql.org/docs/current/tutorial-populate.html)

## INSERT

Non-numeric literals use **single quotes** `'...'`. `date` accepts several formats; the tutorial uses unambiguous `YYYY-MM-DD`.

```sql
INSERT INTO weather VALUES ('San Francisco', 46, 50, 0.25, '1994-11-27');
```

`point` input:

```sql
INSERT INTO cities VALUES ('San Francisco', '(-194.0, 53.0)');
```

**Explicit columns** (recommended style): order can differ; omit columns if allowed (e.g. unknown `prcp`):

```sql
INSERT INTO weather (city, temp_lo, temp_hi, prcp, date)
    VALUES ('San Francisco', 43, 57, 0.0, '1994-11-29');

INSERT INTO weather (date, city, temp_hi, temp_lo)
    VALUES ('1994-11-29', 'Hayward', 54, 37);
```

## COPY (bulk load)

Faster for large loads; path is read by the **server** (backend), not the client:

```sql
COPY weather FROM '/home/user/weather.txt';
```

Tab-separated example row file uses `\N` for NULL. See [COPY](https://www.postgresql.org/docs/current/sql-copy.html).

Prev: [03_create_table.md](03_create_table.md) · Next: [05_query.md](05_query.md)

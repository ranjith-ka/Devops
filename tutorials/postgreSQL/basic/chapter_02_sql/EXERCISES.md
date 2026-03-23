# Hands-on practice — Chapter 2 SQL

Work through these in order. Use **`psql`** against a database you own (e.g. `mydb`). Adjust paths if your clone lives elsewhere.

**Reference:** section notes in this folder and the [official tutorial](https://www.postgresql.org/docs/current/tutorial-sql.html).

---

## 0. Connect

```bash
psql -d mydb
```

In `psql`, useful meta-commands:

- `\dt` — list tables  
- `\d weather` — describe table `weather`  
- `\q` — quit  

---

## 1. Reset lesson data

From your shell (repo root or use the full path to `basics.sql`):

```bash
psql -d mydb -f minikube/postgreSQL/chapter_02_sql/basics.sql
```

Or inside `psql` (path = where **your** file is):

```sql
\i minikube/postgreSQL/chapter_02_sql/basics.sql
```

**Check:** `\dt` shows `cities` and `weather`. `SELECT count(*) FROM weather;` should return **3**.

---

## 2. SELECT and expressions

**Goal:** Read all columns, then compute an average temperature per row.

1. List every column from `weather` without using `*`.
2. Same data using `SELECT * FROM weather;`.
3. Add a column that shows `(temp_hi + temp_lo) / 2` and label it `temp_avg` with `AS`.

**Check:** Step 3 returns **3 rows**; `temp_avg` for the first Mumbai row should be **48**.

<details>
<summary>Solution</summary>

```sql
SELECT city, temp_lo, temp_hi, prcp, date FROM weather;

SELECT * FROM weather;

SELECT city, (temp_hi + temp_lo) / 2 AS temp_avg, date FROM weather;
```

</details>

---

## 3. WHERE and ORDER BY

**Goal:** Filter and sort.

1. Rows where `city` is `'Mumbai'`.
2. Rows where `prcp` is greater than `0` (hint: one row).
3. All rows ordered by `date`, then by `city`.

**Check:** Step 2 returns **1 row**. Step 3 always shows **3 rows** in a deterministic order for the same data.

<details>
<summary>Solution</summary>

```sql
SELECT * FROM weather WHERE city = 'Mumbai';

SELECT * FROM weather WHERE prcp > 0;

SELECT * FROM weather ORDER BY date, city;
```

</details>

---

## 4. DISTINCT

**Goal:** Unique city names from `weather`.

1. `SELECT DISTINCT city FROM weather ORDER BY city;`

**Check:** **2** rows: `Mumbai`, `Pune`.

<details>
<summary>Solution</summary>

```sql
SELECT DISTINCT city FROM weather ORDER BY city;
```

</details>

---

## 5. INSERT (extend the data)

**Goal:** Add rows so joins and aggregates are more interesting.

1. Insert a city **`Bengaluru`** with location **`(77.6, 12.9)`** (same pattern as Mumbai in `basics.sql`).
2. Insert a **`weather`** row: city **`Bengaluru`**, `temp_lo` **40**, `temp_hi` **58**, `prcp` **0**, date **`1994-11-30`** (use explicit column list).
3. Insert another **`weather`** row for **`Pune`** on **`1994-11-30`**: low **36**, high **60**, precipitation **0.1**.

**Check:** `SELECT count(*) FROM cities;` → **2**. `SELECT count(*) FROM weather;` → **5**.

<details>
<summary>Solution</summary>

```sql
INSERT INTO cities VALUES ('Bengaluru', '(77.6, 12.9)');

INSERT INTO weather (city, temp_lo, temp_hi, prcp, date)
    VALUES ('Bengaluru', 40, 58, 0.0, '1994-11-30');

INSERT INTO weather (city, temp_lo, temp_hi, prcp, date)
    VALUES ('Pune', 36, 60, 0.1, '1994-11-30');
```

</details>

---

## 6. INNER JOIN

**Goal:** Match `weather` to `cities` by name.

1. Join `weather` and `cities` so that `weather.city = cities.name`. Start with `SELECT *` and an explicit `JOIN ... ON`.
2. Repeat with only: `city`, `temp_lo`, `temp_hi`, `date`, `location`.

**Check:** Inner join returns **4 rows** (Pune’s two weather rows have **no** matching city in `cities` after section 5).

<details>
<summary>Solution</summary>

```sql
SELECT * FROM weather JOIN cities ON weather.city = cities.name;

SELECT w.city, w.temp_lo, w.temp_hi, w.date, c.location
    FROM weather w
    JOIN cities c ON w.city = c.name;
```

</details>

---

## 7. LEFT OUTER JOIN

**Goal:** Keep all weather rows, even without a city.

1. `LEFT OUTER JOIN` `weather` to `cities` on `weather.city = cities.name` with `SELECT *`.
2. Confirm **Pune** rows show **NULL** in `cities` columns.

**Check:** **5 rows** total.

<details>
<summary>Solution</summary>

```sql
SELECT *
    FROM weather w
    LEFT OUTER JOIN cities c ON w.city = c.name;
```

</details>

---

## 8. Aggregates, GROUP BY, HAVING

**Goal:** Summarize by city.

1. Per `city`: `count(*)`, `max(temp_hi)`, `min(temp_lo)`.
2. Same, but only cities where `max(temp_lo) < 45` (use `HAVING`).

**Check:** Step 1 has **3** groups (Bengaluru, Mumbai, Pune). Step 2 leaves **Bengaluru** and **Pune** (Mumbai’s max low is **46**).

<details>
<summary>Solution</summary>

```sql
SELECT city, count(*), max(temp_hi), min(temp_lo)
    FROM weather
    GROUP BY city;

SELECT city, count(*), max(temp_hi), min(temp_lo)
    FROM weather
    GROUP BY city
    HAVING max(temp_lo) < 45;
```

</details>

---

## 9. Subquery with aggregate

**Goal:** Cities that recorded the **global** maximum `temp_hi` in this table.

1. Find `max(temp_hi)` from `weather`.
2. Select `city` (and optionally `date`) for rows where `temp_hi` equals that maximum.

**Check:** At least **Pune** on `1994-11-30` (`temp_hi` **60**). If only one row has that high, one row; if tied, multiple.

<details>
<summary>Solution</summary>

```sql
SELECT max(temp_hi) FROM weather;

SELECT city, date, temp_hi
    FROM weather
    WHERE temp_hi = (SELECT max(temp_hi) FROM weather);
```

</details>

---

## 10. UPDATE

**Goal:** Fix data in place.

1. Subtract **1** from both `temp_lo` and `temp_hi` for all rows where `date = '1994-11-30'`.
2. `SELECT * FROM weather ORDER BY city, date;` and confirm the change.

**Check:** Bengaluru becomes **39** / **57** for that date after the update.

<details>
<summary>Solution</summary>

```sql
UPDATE weather
    SET temp_lo = temp_lo - 1, temp_hi = temp_hi - 1
    WHERE date = '1994-11-30';

SELECT * FROM weather ORDER BY city, date;
```

</details>

---

## 11. DELETE (scoped)

**Goal:** Remove rows safely.

1. Delete only **`weather`** rows where `city = 'Bengaluru'`.
2. Optionally remove **`Bengaluru`** from **`cities`** (no longer referenced by `weather`).

**Check:** No Bengaluru in `weather`; `SELECT * FROM cities;` has only Mumbai if you did step 2.

<details>
<summary>Solution</summary>

```sql
DELETE FROM weather WHERE city = 'Bengaluru';

DELETE FROM cities WHERE name = 'Bengaluru';
```

</details>

---

## 12. Challenge (optional)

**Goal:** Self-join on `weather`.

List pairs of **different** rows (use `w1.city`, `w2.city` or row comparison) where **`w1.temp_hi > w2.temp_hi`** and **`w1.date = w2.date`**. Restrict to **same day** only.

**Hint:** `w1.city <> w2.city` avoids pairing a row with itself.

<details>
<summary>Solution (one approach)</summary>

```sql
SELECT w1.city AS hotter_city, w1.temp_hi, w2.city AS cooler_city, w2.temp_hi, w1.date
    FROM weather w1
    JOIN weather w2
        ON w1.date = w2.date
        AND w1.temp_hi > w2.temp_hi
        AND w1.city <> w2.city;
```

(Exact row count depends on data after your edits.)

</details>

---

## Reset anytime

To return to the **original** three weather rows and one city:

```bash
psql -d mydb -f minikube/postgreSQL/chapter_02_sql/basics.sql
```

Then you can redo exercises from section **2** onward.

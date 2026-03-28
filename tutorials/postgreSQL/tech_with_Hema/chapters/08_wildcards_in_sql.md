## 8. Wildcards in SQL

[▶ Watch this section](https://www.youtube.com/watch?v=btZdKO17xOM&t=2834s) · [⬆ Back to top](#table-of-contents)

Pattern matching with `LIKE`:

```sql
SELECT * FROM employees WHERE name LIKE 'R%';   -- starts with R
SELECT * FROM employees WHERE name LIKE '%vi';   -- ends with vi
SELECT * FROM employees WHERE name LIKE '_a%';   -- second char is 'a'
```

---

---

## Practice

> Write your own queries below.

```sql
-- Your practice queries here

devdb=# select * from github_users where login like '%dev';
 id |    login    | display_name |  company   | created_at
----+-------------+--------------+------------+------------
  1 | ranjith-dev | Ranjith A    | Acme Cloud | 2018-03-12
(1 row)

 mydb=# select * from employees where name like '%R%';
 name | age
------+-----
(0 rows)
```

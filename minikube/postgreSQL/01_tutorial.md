# PostgreSQL Tutorial

https://www.postgresql.org/docs/current/index.html — Official Documentation

**Chapter 2 (SQL language):** [chapter_02_sql/README.md](chapter_02_sql/README.md) — local notes aligned with [Chapter 2. The SQL Language](https://www.postgresql.org/docs/current/tutorial-sql.html).




```mermaid
flowchart TB
  P[postgres — accepts new connections]
  B1[Backend process 1]
  B2[Backend process 2]
  B3[Backend process N]
  C1[Client 1] <--> B1
  C2[Client 2] <--> B2
  CN[Client N] <--> B3
  P -.->|starts| B1
  P -.->|starts| B2
  P -.->|starts| B3
```
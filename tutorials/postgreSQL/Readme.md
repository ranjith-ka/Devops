# PostgreSQL tutorial

## Docker

```bash
cd tutorials/postgreSQL
docker compose up -d
```

First start loads the [GitHub-style practice dataset](practice/github_sample/README.md) into `devdb` automatically. Connect:

```bash
psql postgresql://postgres:postgres@localhost:5432/devdb
```
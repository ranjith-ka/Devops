# 1.1 Installation

Official: [1.1. Installation](https://www.postgresql.org/docs/current/tutorial-install.html)

PostgreSQL might already be on your machine (OS packages or admin install). If so, use your OS or admin docs for **how to connect**.

You can also install it yourself; a normal user can install PostgreSQL in many setups (**root not always required**).

**Self-install:** Follow [Chapter 17. Installation from Source Code](https://www.postgresql.org/docs/current/installation.html), especially **environment variables** (`PATH`, etc.), then continue this tutorial.

**Remote or non-default server:** Set as needed:

- **`PGHOST`** — database server hostname  
- **`PGPORT`** — port (default is often `5432`)

If clients report they **cannot connect**, verify env and server config with your admin or [client authentication](https://www.postgresql.org/docs/current/client-authentication.html) docs.

**This repo:** Optional local server via [Docker Compose](../docker-compose.yaml) — see [README](README.md#postgresql-in-docker-this-repo).

Next: [02_architecture.md](02_architecture.md)

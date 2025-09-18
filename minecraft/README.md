## Minecraft Project

This project runs a local Minecraft Paper server via Docker and mounts a sample "HelloWorld" datapack that greets all players when the server starts.

### Prerequisites
- Docker Desktop (or Docker Engine) installed and running

### Quick start
```bash
docker compose up -d
```

Then connect to `localhost:25565` from your Minecraft client (matching the configured `VERSION`).

### Files & Folders
- `docker-compose.yml`: Orchestrates the Paper server container
- `datapacks/`: Local datapacks mounted into the world at `/data/world/datapacks`
- `data/`: Server data (worlds, logs, config). Ignored by git

### Configuration
Edit `docker-compose.yml` to tweak:
- `VERSION`: Minecraft/Paper version
- `MEMORY`: JVM memory for the server
- `DIFFICULTY`, `ONLINE_MODE`, etc.

### Datapacks
The included `helloworld` datapack auto-loads using the `minecraft:load` function tag and displays a greeting via `tellraw`. Add your own datapacks in the `datapacks/` folder.

### Useful commands
```bash
# Start in detached mode
docker compose up -d

# View logs
docker compose logs -f

# Stop and remove
docker compose down
```



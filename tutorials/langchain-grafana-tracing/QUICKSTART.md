# 🚀 Quick Start Guide

## What You Have

A production-ready LangChain tracing system with:
- ✅ 7 modular Python files (460 lines, clean architecture)
- ✅ Web UI for testing (2 tabs: Q&A and trace comparison)
- ✅ Docker containerization for easy deployment
- ✅ OpenTelemetry integration with Grafana Tempo
- ✅ Complete documentation and examples

## Run with Docker (Recommended)

```bash
# Build and start all services
docker compose build
docker compose up -d

# Open in browser
# Web UI: http://localhost:5000
# Grafana: http://localhost:3000 (admin/admin)
```

## Run Locally (Development)

```bash
# 1. Create virtual environment
python3 -m venv .venv
source .venv/bin/activate

# 2. Install dependencies
pip install -r requirements.txt

# 3. Run CLI
python app.py "What is LangChain?"
python app.py --compare-trace-a trace-a --compare-trace-b trace-b

# 4. Run Web UI (requires OTLP endpoint running)
python ui.py
# Open: http://localhost:5000
```

## Project Structure

```
├── Core Modules (460 lines total)
│   ├── config.py                    # Configuration
│   ├── trace_data.py               # Data models
│   ├── tracing.py                  # OpenTelemetry setup
│   ├── chain.py                    # LangChain pipeline
│   ├── trace_analyzer.py           # Trace comparison
│   ├── tempo.py                    # Demo data
│   └── app.py                      # CLI orchestration
│
├── Web UI
│   ├── ui.py                       # Flask web server
│   └── templates/
│       ├── base.html               # Base template
│       └── index.html              # Main UI
│
├── Docker
│   ├── Dockerfile                  # Container image
│   ├── docker-compose.yaml         # Stack definition
│   └── otel-collector-config.yaml  # OTLP config
│
├── Documentation
│   ├── ARCHITECTURE.md             # Module architecture
│   ├── DOCKER.md                   # Deployment guide
│   ├── QUICKSTART.md               # This file
│   └── docs/
│       ├── langgraph-architecture.md
│       ├── production-langgraph-architecture.md
│       └── trace-comparison-architecture.md
│
└── Configuration
    ├── requirements.txt            # Python dependencies
    └── README.md                   # Project info
```

## Web UI Features

### Tab 1: Question Answering
- Ask any question
- Get instant answer with trace ID
- View trace details in Grafana

### Tab 2: Compare Traces
- Compare two trace IDs
- See root cause analysis
- View latency differences
- Get recommendations

## API Endpoints (Programmatic Access)

```bash
# Health check
curl http://localhost:5000/health

# Ask question
curl -X POST http://localhost:5000/api/question \
  -H "Content-Type: application/json" \
  -d '{"question": "What is LangChain?"}'

# Compare traces
curl -X POST http://localhost:5000/api/compare \
  -H "Content-Type: application/json" \
  -d '{
    "trace_a": "FIRST_32_CHARACTER_TRACE_ID",
    "trace_b": "SECOND_32_CHARACTER_TRACE_ID"
  }'
```

## Testing Data

Demo traces available for immediate testing:
- Use trace IDs returned by successful question-answering requests.
- Both traces must still be available within Tempo's configured retention window.
- `trace-c` - Trace with validation error

## Key Features

### Single Responsibility Architecture
- Each module has ONE clear purpose
- No circular dependencies
- Easy to test and extend
- ~65 lines per module

### Production Ready
- Health checks built-in
- Docker best practices
- Clean error handling
- Logging and tracing integrated

### Observable
- Full OpenTelemetry instrumentation
- Traces sent to Grafana Tempo
- Visualize with Grafana
- Compare and analyze traces

## Next Steps

1. **Run Docker**: `docker compose up -d`
2. **Open Web UI**: http://localhost:5000
3. **Try Q&A**: Ask a question in Tab 1
4. **Compare Traces**: Use Tab 2 with demo IDs
5. **View in Grafana**: http://localhost:3000

## Documentation

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Module design
- [DOCKER.md](./DOCKER.md) - Deployment guide
- [docs/langgraph-architecture.md](./docs/langgraph-architecture.md) - System design

## Troubleshooting

**Web UI not working?**
```bash
docker compose logs langchain-tracing-ui
docker compose restart langchain-tracing-ui
```

**Traces not in Grafana?**
```bash
docker compose logs otel-collector
docker compose logs tempo
```

**Want to rebuild?**
```bash
docker compose down -v
docker compose build --no-cache
docker compose up -d
```

## Support

Check the DOCKER.md file for detailed troubleshooting and production deployment guidance.

---

**Ready to go\!** 🎉 Your tracing system is now containerized and ready for testing.

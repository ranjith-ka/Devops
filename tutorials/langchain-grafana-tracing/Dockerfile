FROM python:3.10-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY config.py .
COPY trace_data.py .
COPY tracing.py .
COPY chain.py .
COPY trace_analyzer.py .
COPY tempo.py .
COPY loki.py .
COPY documentation.py .
COPY graph.py .
COPY ui.py .
COPY app.py .
COPY templates/ ./templates/
COPY docs/ ./docs/

# Environment variables
ENV PYTHONUNBUFFERED=1
ENV SERVICE_NAME=langchain-tracing
ENV OTLP_ENDPOINT=http://otel-collector:4317

# Expose ports
EXPOSE 5000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:5000/health')" || exit 1

# Start the web UI
CMD ["python", "ui.py"]

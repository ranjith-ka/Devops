"""Configuration from environment variables."""

import os
from urllib.parse import urlparse


def normalize_otlp_endpoint(value: str | None) -> str:
    """Normalize OTLP endpoint values for the Python gRPC exporter."""
    if not value:
        return "localhost:4317"

    candidate = value.strip().rstrip("/")
    parsed = urlparse(candidate)
    if parsed.scheme and parsed.netloc:
        host = parsed.hostname or parsed.netloc
        port = parsed.port
        if port:
            return f"{host}:{port}"
        return host

    return candidate.replace("http://", "").replace("https://", "")


class Config:
    """Application configuration."""

    SERVICE_NAME = os.getenv("OTEL_SERVICE_NAME") or os.getenv("SERVICE_NAME", "langchain-tutorial")
    SERVICE_VERSION = os.getenv("SERVICE_VERSION", "0.1.0")
    ENVIRONMENT = os.getenv("DEPLOYMENT_ENVIRONMENT") or os.getenv("ENVIRONMENT", "local")
    OTLP_ENDPOINT = normalize_otlp_endpoint(
        os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT") or os.getenv("OTLP_ENDPOINT", "http://localhost:4317")
    )
    TEMPO_ENDPOINT = os.getenv("TEMPO_ENDPOINT", "http://localhost:3200").rstrip("/")
    LOKI_ENDPOINT = os.getenv("LOKI_ENDPOINT", "http://localhost:3100").rstrip("/")
    CHECKPOINT_DB = os.getenv("CHECKPOINT_DB", "./data/checkpoints.sqlite")
    DOCUMENTATION_PATH = os.getenv("DOCUMENTATION_PATH", "./docs")
    DOCUMENT_TOP_K = int(os.getenv("DOCUMENT_TOP_K", "3"))
    TRACE_CONTENT = os.getenv("TRACE_CONTENT", "false").lower() == "true"

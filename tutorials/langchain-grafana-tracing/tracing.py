"""OpenTelemetry tracing setup and instrumentation."""

from __future__ import annotations

import time
from collections.abc import Callable
from typing import Any

from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace import Status, StatusCode

from config import Config

TRACER_PROVIDER = None


def configure_tracing() -> TracerProvider:
    """Configure OpenTelemetry tracing with OTLP exporter."""
    global TRACER_PROVIDER
    if TRACER_PROVIDER is not None:
        return TRACER_PROVIDER

    resource = Resource.create(
        {
            "service.name": Config.SERVICE_NAME,
            "service.version": Config.SERVICE_VERSION,
            "deployment.environment.name": Config.ENVIRONMENT,
        }
    )
    provider = TracerProvider(resource=resource)
    exporter = OTLPSpanExporter(endpoint=Config.OTLP_ENDPOINT, insecure=True)
    provider.add_span_processor(BatchSpanProcessor(exporter, schedule_delay_millis=1000))
    trace.set_tracer_provider(provider)
    TRACER_PROVIDER = provider
    return provider


def get_tracer(name: str = "langchain.tutorial"):
    """Return the active tracer, creating the real provider if needed."""
    provider = trace.get_tracer_provider()
    if type(provider).__name__ == "NoOpTracerProvider":
        configure_tracing()
    return trace.get_tracer(name)


def traced_step(name: str, operation: str, function: Callable[[Any], Any]):
    """Decorator to wrap a function in an OpenTelemetry span."""
    from langchain_core.runnables import RunnableLambda

    def invoke(value: Any) -> Any:
        started = time.perf_counter()
        tracer = get_tracer()
        with tracer.start_as_current_span(name) as span:
            span.set_attribute("gen_ai.operation.name", operation)
            span.set_attribute("langchain.step.name", name)
            try:
                result = function(value)
                duration_ms = (time.perf_counter() - started) * 1000
                span.set_attribute("step.duration_ms", duration_ms)
                return result
            except Exception as error:
                span.record_exception(error)
                span.set_status(Status(StatusCode.ERROR, str(error)))
                raise

    return RunnableLambda(invoke)

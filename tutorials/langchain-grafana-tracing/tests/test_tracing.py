import unittest

from opentelemetry import trace as otel_trace

from tracing import configure_tracing, get_tracer


class TraceSetupTest(unittest.TestCase):
    def test_configure_tracing_creates_tracer(self):
        provider = configure_tracing()
        tracer = get_tracer("demo")
        with tracer.start_as_current_span("demo-span") as span:
            self.assertIsNotNone(span)
        provider.force_flush()
        provider.shutdown()

    def test_get_tracer_initializes_real_provider(self):
        global_provider = otel_trace.get_tracer_provider()
        provider_name = type(global_provider).__name__
        self.assertNotEqual(provider_name, "NoOpTracerProvider")
        tracer = get_tracer("demo")
        with tracer.start_as_current_span("lazy-span") as span:
            self.assertIsNotNone(span)
        trace_provider = otel_trace.get_tracer_provider()
        self.assertNotEqual(type(trace_provider).__name__, "NoOpTracerProvider")


if __name__ == "__main__":
    unittest.main()

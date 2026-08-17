#!/usr/bin/env python3
"""Local demo runner - Test the tracing system without external dependencies."""

import json
import sys

# Simulate the core logic without external imports
class TraceSpan:
    def __init__(self, name, duration_ms, status="ok", metadata=None):
        self.name = name
        self.duration_ms = duration_ms
        self.status = status
        self.metadata = metadata or {}

class TraceData:
    def __init__(self, trace_id, root_duration_ms, root_status, spans):
        self.trace_id = trace_id
        self.root_duration_ms = root_duration_ms
        self.root_status = root_status
        self.spans = spans

def create_demo_trace(trace_id, model_delay=180, tool_delay=80, retrieval_delay=120):
    """Create a demo trace for testing."""
    spans = [
        TraceSpan("prompt.render", 15, "ok"),
        TraceSpan("model.generate", model_delay, "ok"),
        TraceSpan("tool.execution", tool_delay, "ok"),
        TraceSpan("retrieval", retrieval_delay, "ok"),
        TraceSpan("output.parse", 20, "ok"),
    ]
    
    total_ms = sum(s.duration_ms for s in spans)
    return TraceData(trace_id, total_ms, "ok", spans)

def compare_traces(trace_a, trace_b):
    """Compare two traces and identify differences."""
    print("\n" + "="*70)
    print(f"📊 TRACE COMPARISON: {trace_a.trace_id} vs {trace_b.trace_id}")
    print("="*70)
    
    print(f"\n📈 Total Duration:")
    print(f"  Trace A: {trace_a.root_duration_ms}ms")
    print(f"  Trace B: {trace_b.root_duration_ms}ms")
    print(f"  Difference: {trace_b.root_duration_ms - trace_a.root_duration_ms}ms")
    
    print(f"\n📍 Span-by-Span Comparison:")
    for span_a, span_b in zip(trace_a.spans, trace_b.spans):
        delta = span_b.duration_ms - span_a.duration_ms
        indicator = "⚠️ " if delta > 0 else "✓ "
        print(f"  {indicator}{span_a.name}:")
        print(f"     A: {span_a.duration_ms}ms  |  B: {span_b.duration_ms}ms  |  Δ: {delta:+d}ms")
    
    # Identify slowest span
    slowest = max(zip(trace_b.spans, trace_a.spans), 
                  key=lambda x: x[0].duration_ms - x[1].duration_ms)
    print(f"\n🎯 Root Cause: {slowest[0].name.upper()} is {slowest[0].duration_ms - slowest[1].duration_ms}ms slower")
    
    if "model" in slowest[0].name:
        print("💡 Recommendation: Check model provider load or token count")
    elif "retrieval" in slowest[0].name:
        print("💡 Recommendation: Optimize database queries or indexing")
    else:
        print("💡 Recommendation: Profile this span for bottlenecks")

def main():
    """Main demo runner."""
    print("\n" + "🚀 "*20)
    print("\n   LangChain Tracing System - Local Demo")
    print("\n" + "🚀 "*20 + "\n")
    
    # Demo 1: Create and display a trace
    print("\n1️⃣  DEMO: Create Baseline Trace")
    print("-" * 70)
    trace_a = create_demo_trace("trace-a", model_delay=100, tool_delay=50, retrieval_delay=80)
    print(f"✓ Created trace: {trace_a.trace_id}")
    print(f"  Total duration: {trace_a.root_duration_ms}ms")
    print(f"  Spans: {len(trace_a.spans)}")
    for span in trace_a.spans:
        print(f"    - {span.name}: {span.duration_ms}ms")
    
    # Demo 2: Create a trace with slower model
    print("\n\n2️⃣  DEMO: Create Trace with Slow Model")
    print("-" * 70)
    trace_b = create_demo_trace("trace-b", model_delay=280, tool_delay=50, retrieval_delay=80)
    print(f"✓ Created trace: {trace_b.trace_id}")
    print(f"  Total duration: {trace_b.root_duration_ms}ms")
    print(f"  Spans: {len(trace_b.spans)}")
    for span in trace_b.spans:
        print(f"    - {span.name}: {span.duration_ms}ms")
    
    # Demo 3: Compare traces
    print("\n\n3️⃣  DEMO: Analyze Differences")
    compare_traces(trace_a, trace_b)
    
    # Demo 4: Create trace with retrieval issue
    print("\n\n4️⃣  DEMO: Create Trace with Retrieval Latency")
    print("-" * 70)
    trace_c = create_demo_trace("trace-c", model_delay=100, tool_delay=50, retrieval_delay=350)
    print(f"✓ Created trace: {trace_c.trace_id}")
    print(f"  Total duration: {trace_c.root_duration_ms}ms")
    print(f"  Spans: {len(trace_c.spans)}")
    for span in trace_c.spans:
        print(f"    - {span.name}: {span.duration_ms}ms")
    
    # Demo 5: Compare a vs c
    print("\n\n5️⃣  DEMO: Identify Retrieval Issue")
    compare_traces(trace_a, trace_c)
    
    # Summary
    print("\n" + "="*70)
    print("✅ LOCAL DEMO COMPLETE")
    print("="*70)
    print("""
What you just saw:

1. ✓ Trace Creation: Generated synthetic traces with realistic timings
2. ✓ Trace Comparison: Compared two traces to identify performance deltas
3. ✓ Root Cause Analysis: Identified which span was the bottleneck
4. ✓ Recommendations: Provided actionable next steps

For the full Web UI experience:

  A. Run with Docker (Recommended):
     docker compose build
     docker compose up -d
     # Then open: http://localhost:5000

  B. Run Locally (requires network):
     python3 -m venv .venv
     source .venv/bin/activate
     pip install -r requirements.txt
     python ui.py
     # Then open: http://localhost:5000

For Grafana traces:
  • Open http://localhost:3000 (admin/admin)
  • Go to Explore → Tempo
  • Paste trace IDs to view full trace hierarchies
""")
    print("🎉 Ready to scale! 🎉\n")

if __name__ == "__main__":
    main()

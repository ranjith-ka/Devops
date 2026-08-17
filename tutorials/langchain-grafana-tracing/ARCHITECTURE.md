# Single Responsibility Module Architecture

Each file has exactly one responsibility:

## 1. **config.py** (14 lines)
Responsibility: Environment configuration
- Centralized config with environment variable loading
- No dependencies on other modules

## 2. **trace_data.py** (94 lines)
Responsibility: Data models
- TraceSpan: represents a single span
- TraceData: represents a complete trace
- SpanComparison: represents span-level comparison
- TraceComparisonResult: represents the full comparison result
- No business logic, only data representation

## 3. **tracing.py** (53 lines)
Responsibility: OpenTelemetry infrastructure
- configure_tracing(): setup OTLP exporter
- TRACER: global tracer instance
- traced_step(): decorator for instrumentation
- Depends only on: config

## 4. **chain.py** (47 lines)
Responsibility: LangChain model and pipeline
- fake_model(): deterministic model implementation
- build_chain(): construct the prompt→model→parser pipeline
- Depends only on: config, tracing

## 5. **trace_analyzer.py** (68 lines)
Responsibility: Trace comparison and root cause diagnosis
- diagnose_root_cause(): identify the most likely cause of differences
- compare_traces(): compare two traces and generate report
- Depends only on: trace_data

## 6. **tempo.py** (76 lines)
Responsibility: Data fetching from Tempo/demo
- make_demo_trace(): generate synthetic trace data for testing
- fetch_trace_from_tempo(): fetch trace (demo implementation uses synthetic data)
- Depends only on: trace_data

## 7. **app.py** (57 lines)
Responsibility: CLI orchestration
- run_question_answering(): orchestrate QA workflow
- run_trace_comparison(): orchestrate comparison workflow
- main(): argument parsing and delegation
- Depends on: everything (it's the entry point)

## Dependency Graph

```
config.py
├── tracing.py
├── chain.py
└── trace_analyzer.py ← trace_data.py
    └── tempo.py ← trace_data.py
        └── app.py (orchestrator)
```

No circular dependencies.
Total lines: 460 lines (production code)

## Testing Individual Modules

```bash
# All modules compile cleanly
for f in config.py trace_data.py tracing.py chain.py trace_analyzer.py tempo.py app.py; do
  python3 -m py_compile "$f"
done
```

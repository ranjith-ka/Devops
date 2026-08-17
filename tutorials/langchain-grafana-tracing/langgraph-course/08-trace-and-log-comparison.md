# Lesson 8: Trace and log comparison

## Learning objective

Understand how raw Tempo OTLP JSON becomes a shared-scale waterfall and how Loki
events add diagnostic context.

## 1. Tempo normalization

`tempo.py` validates a 32-character trace ID, calls Tempo, combines all returned
batches, identifies the root span, and converts nanoseconds into milliseconds.
Each child retains its start offset, duration, status, and metadata.

## 2. Span matching

`trace_analyzer.py` orders spans by start time and assigns occurrence keys:

```text
tool
tool [2]
tool [3]
```

This matters when future graphs retry a node or execute the same tool repeatedly.
Matching only by name would silently compare the wrong calls.

## 3. Diagnosis threshold

A regression must be at least 25 ms and at least 10% of Trace A’s duration for
that span. This filters tiny runtime jitter. The UI still shows every raw delta.

## 4. Combined report

The comparison endpoint fetches both Tempo traces and both Loki event streams.
The UI uses one shared maximum duration so horizontal bar positions are directly
comparable.

## Exercise

Compare two requests on the same thread and answer:

1. Which root request was faster?
2. Which node contributed the largest absolute difference?
3. Did both traces retrieve the same document count?
4. Did memory contain a different number of messages?

## Next lesson

Lesson 9 provides repeatable recipes for adding nodes, routing, tools, retries,
and validation.

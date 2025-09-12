# emapy - Logistics AI Demo (Docker Compose)

This is a demo scaffold:
- React frontend (UI stub)
- Go backend API (prioritization, dedupe, HS suggestions, bulk updates, assignment)
- Postgres DB (demo)
- ML microservice stubs: NLP, CV, TimeSeries, Optimizer
- Docker Compose for local demo

Quick start:
1. Review files.
2. chmod +x demo-run.sh
3. ./demo-run.sh
4. Open http://localhost


Mapped features & implementation notes (what to build and how)

Order Prioritization (automation + manual override) — NLP + rules

Auto-tag incoming orders using metadata + NLP on email/IDE text: extract ETA, SLA, customer tier, “hypercare” flags.

Rules engine for priority scoring (weight ETA, SLA, customer expectation, hypercare).

UI: prioritized inbox with manual override and audit trail.

Tech: Go API (or existing API) calls a small NLP microservice (HuggingFace/DistilBERT) for tag extraction + a lightweight rules engine (JSON rules, Drools or simple evaluation).

Assignment of Orders — Marketplace / self-assign + suggestions

Provide CPA-facing “suggested assignment” list based on skills, workload, declaration types.

Integrate with availability/shift data (calendar).

UI: “claim” button; show reason why suggested.

HS Code Classification & Declaration Preparation — Hybrid ML + tool integration

Provide ML-assisted HS code suggestions (model + confidence) and link to United Wizard for verification.

Allow batch/inline editing; enable bulk updates across items and cases.

Tech: classification microservice (fine-tune a classifier on historical HS mappings), fallback to United Wizard integration API.

Duplicate Case Detection — Deterministic + ML fuzzy matching

Use thread headers, sender email fingerprints, subject similarity, and dedupe heuristics.

Maintain a case-hash index and similarity threshold. Mark probable duplicates for review and provide dedupe suggestions.

Bulk Update Enhancements

Bulk UI for HS, Y codes, preferences, supplementary units. Rollback capability.

Server-side validation and idempotent operations.

Registry Enrichment (Transport Unit Registry)

Integrate external registry API to enrich consignment headers (auto-populate fields).

Cache registry results and allow manual override.

Milestone & Status Tracking

Support multiple declarations per case, each with own milestone timeline; display in dashboards and export to Excel.

UI & Accessibility Improvements

Priority flags, email recall UX (soft-delete / undo), notification improvements, scroll-bar/accessibility CSS changes.

Badge Selection Optimization

Context-driven dropdown (port + goods location) — server-provided lookups to keep client small.

Deployment Timing & Ops

Blue/green or canary deployment windows to avoid border disruptions, maintenance windows UI, feature flags for toggles.

Prioritized backlog (P0 = must-have for MVP demo; P1 = important; P2 = nice-to-have)

P0

Order Prioritization auto-tagging + priority inbox (NLP+rules).

Duplicate case detection (subject+thread heuristics).

HS-code suggestion assistant (single-item UI + confidence).

Bulk update UI (HS/Y codes + preferences) — basic rollback.

Basic CPA assignment UI (self-assign + suggested).

P1

Registry enrichment microservice + caching.

Multi-declaration milestone tracking + Excel export.

Assignment suggestion using workload/skills.

UI polish: high-priority flags, notification logic improvements.

P2

Full ML dedupe (fuzzy clustering), automated preference deployment, advanced badge selection logic, email recall (with retention/undo service), advanced accessibility (contrast/themes).

Acceptance criteria / DoD (per P0 item)

Auto-tagging: given a set of test emails, system auto-tags with ≥ X% precision on defined labels (X = measurable baseline). UI displays suggested tag and user can override. Audit log records overrides.

Duplicate detection: system flags duplicates for 90% of known duplicate test threads (measured on sample), and no more than Y% false positives.

HS-code assistant: suggests top-3 HS codes with confidence score; manual override persists chosen code and logs decision.

Bulk update: UI allows selecting N items and apply HS code change; operation is atomic per-case or transactional per-item; rollback button reverts last bulk op.

Data & training needs

Historical case data (anonymized): subject, body, thread id, assigned labels, HS mappings — needed for HS classifier and dedupe model.

Sample emails and attachments for NLP extraction.

Telemetry for performance of bulk updates (size/performance).

Integration credentials for United Wizard and Transport Unit Registry (or simulation/test endpoints).

Monitoring, alerts & SLOs

Metrics: priority_inference_latency, duplicate_detection_rate, bulk_update_duration, hs_suggestion_accuracy (if CI pipeline runs validations).

Dashboards: prioritized inbox health, duplicate rate over time, number of manual overrides (measure model trust), bulk update throughput.

Alerts:

High duplicate rate increase (sudden spike) → Ops investigate.

Bulk-update failures > 1% → Abort and notify.

Registry API failure → Degrade gracefully and create incident.

Deployment & release checklist (DevOps)

Build images (frontend, go-api, nlp, registry-enricher, etc.) with multi-stage Dockerfiles.

CI pipeline: run unit tests, static analysis, container image scan (Trivy).

Feature flags and canary rollout: deploy auto-tagging behind a flag; start with 10% traffic.

DB migrations backup and rollback plan (schema changes for multi-declaration).

Pre-warm models and have fallbacks (stubs) for demo.

Observability: Prometheus scrape + Grafana dashboard + structured logs shipped to central store.

Security & compliance

Anonymize PII used for model training and demos. Redact emails when sharing datasets.

For model components check licenses (YOLOv5/others) and include legal.txt.

Apply RBAC on assignment operations and audit logs for manual overrides.

Demo / customer scenarios (3 quick scripts)

Prioritization demo: ingest sample batch of orders; show auto-tags + priority score; demonstrate manual override & audit trail.

Duplicate detection demo: feed overlapping email threads and show dedupe suggestions and how to merge or close duplicate.

HS + bulk update demo: show HS suggestion for a single ad-hoc case, then perform a bulk HS update across 20 items and rollback.

Risks & mitigation

False positives in dedupe/HS models → show confidence & require human-in-loop for initial rollout.

Registry API outages → cache last-good and expose manual enrichment flows.

Regulatory / license issues → prefer permissively licensed models for production; keep demo stubs.

Ready-to-use artifacts I can produce right now

OpenAPI spec for the Go API endpoints (cases, detect, bulk-update).

Prototype Docker Compose with NLP stub + dedupe rule-engine + frontend demo route.

Short technical slide (1 page) mapping business value to features & KPIs for customer pitch.

Acceptance test suite (curl + sample payloads) for demo flows.

Pick any one artifact from the list above and I’ll generate it immediately.
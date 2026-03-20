# Writing prompts for AI agents

This note describes a simple structure for prompts so agents behave predictably: clear identity, bounded work, enough background, explicit thinking style, when to stop, and what to deliver.

---

## 1. Role

**What it is:** Who the model should *be*—expertise, tone, and constraints on behavior.

**Why it matters:** Sets default assumptions (e.g. “senior SRE” vs “tutorial writer”) and reduces generic or off-brand answers.

**How to write it:**

- Name the persona in one line (role + seniority or domain).
- Add 1–3 hard rules: what you always do, what you never do.
- Optional: audience (e.g. “for engineers who know Kubernetes”).

**Example:**

```text
You are a senior platform engineer. You prefer concrete commands and manifests over vague advice. You never suggest disabling security controls without naming the tradeoff.
```

---

## 2. Task

**What it is:** The specific outcome you want—verbs, scope, and success criteria in plain language.

**Why it matters:** Without a task, the model fills space with overview instead of doing the work.

**How to write it:**

- Start with an action verb: *draft*, *refactor*, *compare*, *debug*, *list steps*, *review*.
- Bound the scope: one feature, one file, one decision—not “everything about X.”
- If there are choices, say which default to use when ambiguous.

**Example:**

```text
Task: Propose a minimal Helm values change to add a readiness probe for the API deployment. Do not change unrelated charts.
```

---

## 3. Context

**What it is:** Facts the model cannot infer: environment, versions, constraints, prior decisions, and relevant excerpts (logs, configs, error messages).

**Why it matters:** Grounds answers in *your* situation instead of generic patterns.

**How to write it:**

- Prefer paste or cite: file paths, snippets, command output, API errors.
- State versions (runtime, cloud, framework) when behavior depends on them.
- Separate **facts** from **hypotheses** so the model does not treat guesses as truth.

**Example:**

```text
Context:
- Cluster: minikube, Kubernetes 1.28
- Error: Pod CrashLoopBackOff; last log line: "connection refused localhost:5432"
- Postgres runs in another namespace; service name `postgres.data.svc.cluster.local`
```

---

## 4. Reasoning

**What it is:** How you want the model to think *before* or *while* answering: chain-of-thought style, tradeoff analysis, or “assume nothing—verify from context.”

**Why it matters:** Improves correctness on multi-step problems and forces explicit assumptions.

**How to write it:**

- Say whether reasoning should be **visible** (show steps) or **internal** (brief answer only).
- Ask for **assumptions** and **unknowns** when the context is incomplete.
- For reviews: “list issues by severity, then recommendations.”

**Example:**

```text
Reasoning:
- Briefly list assumptions, then give the fix.
- If the root cause is unclear from logs, state what extra command output you need.
```

---

## 5. Stop conditions

**What it is:** When the model should **stop** and hand back: limits on length, depth, file count, or “stop after first working solution.”

**Why it matters:** Prevents endless elaboration, duplicate suggestions, or scope creep.

**How to write it:**

- Cap output: max sections, max bullet depth, or “one page.”
- Cap exploration: “do not propose alternatives unless the primary approach fails.”
- Escalation: “if a secret or credential is required, stop and list what to supply.”

**Example:**

```text
Stop conditions:
- Deliver at most one recommended approach plus a short fallback.
- Stop after the minimal diff; do not rewrite unrelated modules.
- If you cannot verify from the given context, say what is missing instead of guessing.
```

---

## 6. Outputs

**What it is:** The **shape** of the answer: sections, format (Markdown, JSON, table), code fences, and what to omit.

**Why it matters:** Makes results easy to paste into tickets, PRs, or runbooks.

**How to write it:**

- List required headings or fields in order.
- Specify code vs prose: “only changed files in fenced blocks with paths.”
- Say what **not** to include: apologies, filler, or repeated context.

**Example:**

```text
Outputs:
1. Summary (3 bullets max)
2. Steps or diff (numbered)
3. Risks / rollback (short)
Use Markdown. No preamble like "Sure, I can help."
```

---

## Full prompt skeleton

Copy and fill in the brackets:

```text
Role:
[persona + non-negotiable behaviors]

Task:
[one clear outcome + scope boundaries]

Context:
[facts, versions, snippets, links or paths]

Reasoning:
[how to think; assumptions; what to do if blocked]

Stop conditions:
[limits; when to stop; no scope creep]

Outputs:
[structure, format, exclusions]
```

---

## Quick checklist

| Piece            | Question to answer                          |
|-----------------|----------------------------------------------|
| **Role**        | Who is speaking, and what rules bind them?   |
| **Task**        | What single thing must be done?              |
| **Context**     | What must the model know that isn’t generic? |
| **Reasoning**   | Show work or stay terse? What if unsure?     |
| **Stop conditions** | When is “enough” enough?                 |
| **Outputs**     | What does “done” look like on the page?      |

Small, explicit prompts with all six elements usually outperform long unstructured requests.

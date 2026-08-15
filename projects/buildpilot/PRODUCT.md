# BuildPilot product definition

## Product promise

BuildPilot answers one question for a coding agent or engineering team:

> Is this AI-generated change safe to merge?

It accepts an immutable base commit plus an optional patch, performs isolated
verification in a customer-controlled Kubernetes cluster, and returns concise,
machine-readable diagnostics. BuildKit is one execution component; remote
verification is the product.

## Initial customer

The first supported customer is a team that uses .NET, GitHub, Kubernetes, and
AI coding agents, and needs builds to remain inside controlled infrastructure.
The initial buyer is a platform-engineering or developer-productivity team.

BuildPilot does not initially target every language, source host, or CI system.

## Sellable workflow

```text
AI agent creates a patch
          |
          v
POST /v1/verifications (repository, base commit, patch, mode)
          |
          v
BuildPilot creates an isolated remote workspace
          |
          v
checkout -> restore -> compile -> affected unit tests
                                   |
                         +---------+---------+
                         |                   |
                    fast result        full verification
                                             |
                              coverage -> Sonar quality gate
                                             |
                              BuildKit -> image digest
                                             |
                              temporary deployment/tests
          |
          v
Structured diagnostics and evidence are returned to the agent
```

`fast` mode optimizes the edit-test loop. It compiles affected projects and
runs relevant unit tests. `full` mode is the merge/release checkpoint and adds
complete tests, coverage, Sonar analysis, container building, and integration
tests.

New requests for the same agent, repository, and base revision supersede older
queued requests. Infrastructure failures may be retried; compilation errors,
test failures, and failed quality gates must not be retried as infrastructure
failures.

## API outcome

The primary output is structured evidence, not a successful Pod or a raw log.

```json
{
  "verification_id": "ver-1042",
  "base_commit": "a81c92f",
  "mode": "fast",
  "status": "failed",
  "stage": "unit_tests",
  "classification": "code_failure",
  "diagnostics": [
    {
      "file": "PaymentServiceTests.cs",
      "line": 74,
      "code": "TEST_ASSERTION",
      "message": "Expected BadRequest, received Accepted"
    }
  ],
  "retryable": false,
  "report_url": "https://buildpilot.example/verifications/ver-1042"
}
```

A successful full verification additionally returns the immutable OCI image
digest, test and coverage summaries, Sonar quality-gate result, and links to
logs, reports, SBOM, and provenance.

## Product boundary

BuildPilot owns the capabilities that distinguish the product:

- patch-based verification requests for coding agents;
- affected-project and affected-test selection;
- fast and full feedback policies;
- cancellation and superseding of obsolete work;
- normalized compiler, test, analysis, and infrastructure diagnostics;
- durable verification history and agent-friendly APIs.

BuildPilot reuses established infrastructure:

- Kubernetes Jobs initially, with Tekton or Argo considered when a workflow
  engine is required;
- BuildKit for OCI image construction;
- SonarQube or SonarCloud for static-analysis quality gates;
- PostgreSQL for durable state;
- an S3-compatible store for logs and reports;
- an existing OCI registry and secret manager.

It is not a replacement for Kubernetes, BuildKit, a registry, Sonar, or a
general-purpose CI workflow engine.

## Runtime architecture

```text
AI agent / GitHub check / CLI
              |
              v
   BuildPilot API and controller
   - authentication and policy
   - durable state machine
   - queue, scheduling, cancellation
   - result normalization
              ^
              | outbound authenticated connection
              v
      customer cluster agent
   - creates isolated Jobs/namespaces
   - streams status and bounded logs
   - uploads reports and evidence
              |
       +------+------+---------+
       |             |         |
  .NET runner    Sonar scan  BuildKit
                             rootless
```

The control plane does not receive customer kubeconfig or registry secrets.
The cluster agent owns namespace-scoped credentials and initiates outbound
communication.

## Verification lifecycle

```text
queued -> provisioning -> restoring -> compiling -> testing
                                                 |         |
                                                 | fast    | full
                                                 v         v
                                           succeeded   analyzing
                                                          |
                                                     image_building
                                                          |
                                                   integration_testing
                                                          |
                                                     succeeded

any active state -> failed | cancelled | superseded | timed_out
```

Each transition is persisted before external effects are acknowledged. The
controller must reconcile after a restart without duplicating a completed
stage. Kubernetes resources are execution state, not the system of record.

## Packaging and delivery

Production delivery is a versioned Helm chart plus signed OCI images:

```text
buildpilot chart
|- control-plane/controller image
|- cluster-agent image
|- dotnet-runner images by supported SDK version
|- CRDs, RBAC, quotas, network policies, and pipeline templates
`- migration and smoke-test Jobs
```

Installation should converge on:

```bash
helm install buildpilot \
  oci://registry.example.com/buildpilot/charts/buildpilot \
  --namespace buildpilot \
  --create-namespace
```

PostgreSQL, object storage, registry, Git provider, Sonar, and secret manager
are configurable integrations rather than mandatory bundled replacements.

A Kind-based bundle may install local substitutes for evaluation and demos. A
prebuilt Kind node image is not the production distribution because customers
will run EKS, AKS, GKE, OpenShift, or another existing Kubernetes platform.

An air-gapped distribution must include an image lock file, Helm chart, all
required images, SBOMs, signatures, and documented upgrade/rollback steps.

## Proof-of-value milestone

The first credible demonstration is:

1. An agent submits a .NET patch without pushing a branch.
2. A remote Kubernetes Job applies it to an immutable base commit.
3. Compilation or a targeted test fails.
4. BuildPilot returns a file, line, category, and non-retryable diagnosis.
5. The agent corrects the patch and fast verification passes.
6. Full verification passes its Sonar quality gate and returns a tested image
   digest.

Initial service objectives are a verification start in under 10 seconds, a
cached fast result in under 60 seconds for the sample repository, no lost state
after a controller restart, and complete cleanup of expired execution resources.

## Repository strategy

BuildPilot remains under `projects/buildpilot` while its API and proof of value
are changing quickly. It is already an independent Go module, so this does not
require coupling it to the rest of this repository.

Move it to a dedicated repository when the first of these occurs:

- an external user or team begins installing it;
- it needs independent releases, CI, issues, access control, or ownership;
- the Helm chart and OCI images are published;
- development cadence starts conflicting with this repository.

When moved, preserve Git history with a subtree split or a repository-filtering
tool, and keep the examples here as consumers of released BuildPilot artifacts.

import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function PipelineGeneratorPage() {
  return (
    <AppShell activePath="/pipeline-generator">
      <ScreenSpec
        title="Pipeline Generator"
        subtitle="Generate GitHub Actions, Azure DevOps YAML, Kubernetes manifests, Terraform, and Helm scaffolds from one governed workflow request."
        api={["POST /generate-pipeline", "POST /generate-terraform", "POST /generate-kubernetes", "POST /explain-workflow"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Input form                           | Generated artifact preview                |
| target, template, guardrails         | YAML, diff, notes, approvals              |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Form > policy guardrails > generated artifact", "Preview > validation > approval bar", "Diff viewer and export actions"]}
        journey={["Choose the target platform and template.", "Apply security and approval guardrails.", "Review the generated output before exporting."]}
        stateModel={["Template, target, security policy, and approval mode drive generation.", "Structured output validation runs before export or copy."]}
        states={[
          { label: 'Loading', description: 'Show form skeletons and a preview placeholder while generation initializes.' },
          { label: 'Error', description: 'Display template validation issues or model generation failures with retry.' },
          { label: 'Empty', description: 'Prompt the user to choose a target platform and a starter template.' }
        ]}
        evidence={["org policy", "workflow context", "security guardrails", "target runtime"]}
      />
    </AppShell>
  );
}

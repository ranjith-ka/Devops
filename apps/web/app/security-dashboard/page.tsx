import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function SecurityDashboardPage() {
  return (
    <AppShell activePath="/security-dashboard">
      <ScreenSpec
        title="Security Dashboard"
        subtitle="See security findings, policy violations, dependency alerts, and AI-generated remediation guidance in one place."
        api={["POST /security-review", "POST /chat"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Risk score                           | Findings and policy violations            |
| severity, trend, SLA                 | clusters, owners, remediations            |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Risk score > findings > remediation queue", "Policy view > owner mapping > exceptions", "AI summary and approval controls"]}
        journey={["Open the latest scan results.", "Review the highest-severity findings and blast radius.", "Generate or approve a remediation plan."]}
        stateModel={["Repo, branch, scan type, policy set, and exception state drive the view.", "High-risk issues are pinned to the top until dismissed or fixed."]}
        states={[
          { label: 'Loading', description: 'Skeleton cards for risk, findings, and owner mapping.' },
          { label: 'Error', description: 'Indicate whether scanner ingestion or policy resolution failed.' },
          { label: 'Empty', description: 'Explain there are no unresolved findings for this repository or scan type.' }
        ]}
        evidence={["SAST findings", "dependency alerts", "container scans", "policy exceptions"]}
      />
    </AppShell>
  );
}

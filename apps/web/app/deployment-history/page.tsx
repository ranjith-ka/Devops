import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function DeploymentHistoryPage() {
  return (
    <AppShell activePath="/deployment-history">
      <ScreenSpec
        title="Deployment History"
        subtitle="Review release progression, outcomes, and rollback decisions with the AI-generated summary for each deployment."
        api={["POST /deployment-summary", "POST /root-cause", "POST /chat"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Release list                         | Release detail and status ribbon          |
| env, version, outcome                | events, approvals, rollback status        |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Release list > status ribbon > detail", "Events > approvals > rollback controls", "AI summary and anomaly flags"]}
        journey={["Select a prior release or environment.", "Audit the event trail and outcome.", "Compare successful and failed deployments over time."]}
        stateModel={["Release, environment, approval, and rollback state drive the detail view.", "Risk and post-release summaries are recalculated when outcomes arrive."]}
        states={[
          { label: 'Loading', description: 'Render release cards and event placeholders while history loads.' },
          { label: 'Error', description: 'Show whether release history, event stream, or approvals failed to load.' },
          { label: 'Empty', description: 'Explain there is no deployment history for the selected service or environment.' }
        ]}
        evidence={["release events", "approvals", "rollback decisions", "post-deploy summaries"]}
      />
    </AppShell>
  );
}

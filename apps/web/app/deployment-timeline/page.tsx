import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function DeploymentTimelinePage() {
  return (
    <AppShell activePath="/deployment-timeline">
      <ScreenSpec
        title="Deployment Timeline"
        subtitle="Track a release from merge to production with risk overlays, stage-by-stage events, and rollback decision points."
        api={["POST /deployment-summary", "POST /root-cause", "POST /chat"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Release list                         | Release detail timeline                   |
| status, env, version                 | stage events, checks, risk overlay        |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Release list > timeline > checks", "Risk overlay > confidence notes > rollback action", "Related logs and commit links"]}
        journey={["Pick a release from the history list.", "Scroll through stage checkpoints and AI notes.", "Decide whether to continue, pause, or rollback."]}
        stateModel={["Release, environment, approval, and event stream state drive the timeline.", "Risk is recomputed when new checks, logs, or security results arrive."]}
        states={[
          { label: 'Loading', description: 'Render release cards, stage markers, and a muted event stream skeleton.' },
          { label: 'Error', description: 'Explain which release metadata or event source is unavailable.' },
          { label: 'Empty', description: 'Show that no deployments exist for the selected service or environment.' }
        ]}
        evidence={["merge commits", "deployment events", "approval gates", "risk score history"]}
      />
    </AppShell>
  );
}

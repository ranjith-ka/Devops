import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function PipelineExplorerPage() {
  return (
    <AppShell activePath="/pipeline-explorer">
      <ScreenSpec
        title="Pipeline Explorer"
        subtitle="Drill into a pipeline graph, inspect stages, and correlate failures with logs, artifacts, and recent code changes."
        api={["POST /analyze-pipeline", "POST /root-cause", "POST /chat"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Pipeline tree                        | Selected pipeline detail                  |
| jobs, stages, runners, artifacts     | runs, logs, duration, flaky tests         |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Pipeline tree > selected run > stage detail", "Evidence panel > log lines > test failures", "Recommended fix actions and rerun controls"]}
        journey={["Select a pipeline from the tree.", "Review the failed stage and associated evidence.", "Ask the assistant to rank the likely root cause."]}
        stateModel={["Pipeline, run, stage, artifact, and commit metadata drive the view.", "Timeline cursor controls which logs, metrics, and PR context are loaded."]}
        states={[
          { label: 'Loading', description: 'Tree skeleton, run card skeletons, and streaming log placeholders.' },
          { label: 'Error', description: 'Tell the user whether pipeline sync, artifact retrieval, or log retrieval failed.' },
          { label: 'Empty', description: 'Explain that no runs exist yet for the selected pipeline and suggest a sync action.' }
        ]}
        evidence={["stage timing", "recent commits", "artifact metadata", "raw logs"]}
      />
    </AppShell>
  );
}

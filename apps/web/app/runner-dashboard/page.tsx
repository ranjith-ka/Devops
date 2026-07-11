import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function RunnerDashboardPage() {
  return (
    <AppShell activePath="/runner-dashboard">
      <ScreenSpec
        title="Runner Dashboard"
        subtitle="Observe runner health, queue depth, utilization, and autoscaling pressure across hosted and self-hosted capacity pools."
        api={["POST /optimize", "POST /analyze-pipeline"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Fleet summary                        | Runner capacity and queue detail          |
| online, offline, utilization         | regions, labels, scaling pressure         |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Fleet summary > capacity pools > queue depth", "Region detail > labels > autoscaling advice", "Failure and saturation alerts"]}
        journey={["Open a runner pool by region or label.", "Check saturation and queue delay trends.", "Apply capacity recommendations or investigate failures."]}
        stateModel={["Runner pool, region, labels, and queue state drive capacity views.", "Autoscaling advice changes with workload mix and sustained saturation."]}
        states={[
          { label: 'Loading', description: 'Render capacity cards and queue charts from cached telemetry.' },
          { label: 'Error', description: 'Show when runner telemetry or availability checks are unavailable.' },
          { label: 'Empty', description: 'Tell the user no runner pools exist for the selected tenant or environment.' }
        ]}
        evidence={["queue depth", "job duration", "runner health", "region usage"]}
      />
    </AppShell>
  );
}

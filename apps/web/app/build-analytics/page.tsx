import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function BuildAnalyticsPage() {
  return (
    <AppShell activePath="/build-analytics">
      <ScreenSpec
        title="Build Analytics"
        subtitle="Analyze build durations, flakiness, test stability, and throughput trends across repositories and branches."
        api={["POST /analyze-pipeline", "POST /optimize", "POST /root-cause"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Trend charts                         | Cohorts and optimization suggestions      |
| build time, pass rate                | caches, parallelism, flaky tests          |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Trends > cohorts > build hotspots", "Optimization notes > test clusters > actions", "Alerting and regression thresholds"]}
        journey={["Select a branch or repo cohort.", "Inspect duration and failure trendlines.", "Use AI to isolate a slowdown or flaky test pattern."]}
        stateModel={["Time range, repo, branch, and job type drive the chart queries.", "Regression detection compares the current cohort to historical baselines."]}
        states={[
          { label: 'Loading', description: 'Use chart skeletons and a summary card while metrics stream in.' },
          { label: 'Error', description: 'Explain whether build metrics, trend aggregation, or cohort lookup failed.' },
          { label: 'Empty', description: 'Show that there is not enough build history yet to chart meaningful trends.' }
        ]}
        evidence={["build durations", "pass rate trends", "test clusters", "cache hit rate"]}
      />
    </AppShell>
  );
}

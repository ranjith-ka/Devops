import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

const metrics = [
  { label: 'Pipeline health', value: '97.8%' },
  { label: 'Deploy risk', value: 'Low' },
  { label: 'Median build', value: '8m 41s' },
  { label: 'Flaky tests', value: '12' }
];

const dashboardStates = [
  { label: 'Loading', description: 'Show KPI skeletons, chart placeholders, and a muted activity feed while tenant data streams in.' },
  { label: 'Error', description: 'Surface the failing service and offer retry plus fallback to cached summaries.' },
  { label: 'Empty', description: 'Explain that the org has no synced repositories yet and offer an integration CTA.' }
];

export default function DashboardPage() {
  return (
    <AppShell activePath="/dashboard">
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {metrics.map((metric) => (
            <article key={metric.label} className="rounded-3xl border border-white/8 bg-white/5 p-5 shadow-glow backdrop-blur">
              <div className="text-sm text-slate-400">{metric.label}</div>
              <div className="mt-2 text-3xl font-semibold tracking-tight">{metric.value}</div>
            </article>
          ))}
        </div>
        <ScreenSpec
          title="Dashboard"
          subtitle="An executive summary for delivery health, risk, and AI recommendations across repositories, pipelines, and deployments."
          api={["POST /chat", "POST /deployment-summary", "POST /optimize"]}
          wireframe={`+--------------------------------------------------------------------------------+
| Topbar: org switcher | search | notifications | profile                        |
+----------------------+---------------------------------------------------------+
| Sidebar              | KPI strip                                               |
| - Dashboard          | [Health] [Deploy risk] [Lead time] [MTTR]               |
| - Pipelines          |                                                         |
| - Deployments        | Trend charts | AI insights | recent failures           |
| - Security           |                                                         |
| - Cost               | Activity feed | release map | recommendations           |
+----------------------+---------------------------------------------------------+`}
          hierarchy={["Shell > KPIs > charts > insights > activity", "KPI cards summarize the live delivery posture", "Action rail surfaces open assistant and export actions"]}
          journey={["Open the org and land on the health summary.", "Review AI explanations for changes in risk and throughput.", "Jump into the most urgent failure, deployment, or security issue."]}
          stateModel={["Org, time range, environment, and repo filters drive the dashboard query model.", "AI summaries are cached by tenant and refresh on data freshness or event triggers."]}
          states={dashboardStates}
          evidence={["pipeline runs", "deployment history", "metrics snapshots", "notification stream"]}
        />
      </div>
    </AppShell>
  );
}

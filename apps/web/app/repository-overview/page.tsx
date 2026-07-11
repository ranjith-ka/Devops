import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function RepositoryOverviewPage() {
  return (
    <AppShell activePath="/repository-overview">
      <ScreenSpec
        title="Repository Overview"
        subtitle="View repository-level health, branch activity, commit velocity, and recent AI-generated recommendations in one place."
        api={["POST /analyze-pipeline", "POST /chat", "POST /deployment-summary"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Repo scorecards                      | Activity feed and health trends           |
| freshness, churn, risk               | commits, PRs, deployments, anomalies      |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Repo scorecards > branch trends > activity feed", "Health signals > ownership > next actions", "AI recommendations and issue links"]}
        journey={["Choose a repository from the tenant inventory.", "Review the current health and volatility indicators.", "Open the most relevant pipeline, PR, or deployment issue."]}
        stateModel={["Repository metadata, branch freshness, recent deployment and PR history drive the score.", "Health recency expires if sync is stale or blocked."]}
        states={[
          { label: 'Loading', description: 'Render scorecard skeletons and a condensed activity stream placeholder.' },
          { label: 'Error', description: 'Surface sync failures, permission issues, or missing repository context.' },
          { label: 'Empty', description: 'Tell the user the repository has not synced delivery metadata yet.' }
        ]}
        evidence={["commit history", "PR activity", "pipeline runs", "release summaries"]}
      />
    </AppShell>
  );
}

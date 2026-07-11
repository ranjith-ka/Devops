import Link from 'next/link';
import { OnboardingShell } from '@/components/onboarding-shell';

const stages = [
  { label: 'Commits', progress: '100%' },
  { label: 'Pipelines', progress: '78%' },
  { label: 'PRs', progress: '62%' },
  { label: 'Deployments', progress: '44%' }
];

export default function SyncProgressPage() {
  return (
    <OnboardingShell activePath="/onboarding/sync">
      <div className="grid gap-6 xl:grid-cols-[1fr_0.85fr]">
        <section className="rounded-[2rem] border border-white/8 bg-white/5 p-6 shadow-glow backdrop-blur">
          <div className="text-xs uppercase tracking-[0.28em] text-cyan-200/80">Sync progress</div>
          <h2 className="mt-3 text-3xl font-semibold tracking-tight">Importing evidence for the first AI summary.</h2>
          <div className="mt-6 space-y-4">
            {stages.map((stage) => (
              <div key={stage.label} className="rounded-3xl border border-white/8 bg-black/20 p-5">
                <div className="flex items-center justify-between text-sm text-slate-200">
                  <span>{stage.label}</span>
                  <span>{stage.progress}</span>
                </div>
                <div className="mt-3 h-2 rounded-full bg-white/10">
                  <div className="h-2 rounded-full bg-cyan-300" style={{ width: stage.progress }} />
                </div>
              </div>
            ))}
          </div>
          <div className="mt-6 flex flex-wrap gap-3">
            <Link href="/onboarding/complete" className="rounded-full bg-cyan-300 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-200">
              View first summary
            </Link>
            <Link href="/onboarding/connect" className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
              Change connection
            </Link>
          </div>
        </section>

        <aside className="space-y-4">
          <div className="rounded-3xl border border-white/8 bg-white/5 p-5 backdrop-blur">
            <div className="text-sm font-medium text-white">Import timeline</div>
            <ul className="mt-3 space-y-3 text-sm leading-7 text-slate-300">
              <li>Pull repos and branches</li>
              <li>Index pipelines and logs</li>
              <li>Map PRs and deployment history</li>
              <li>Generate the first AI summary</li>
            </ul>
          </div>
          <div className="rounded-3xl border border-amber-300/20 bg-amber-300/10 p-5 text-sm leading-7 text-amber-50">
            Partial sync failures are recoverable. The user can resume without losing their place.
          </div>
        </aside>
      </div>
    </OnboardingShell>
  );
}

import Link from 'next/link';
import { OnboardingShell } from '@/components/onboarding-shell';

const nextActions = [
  'Open the dashboard',
  'Review pipeline health',
  'Ask the AI assistant about the repo',
  'Generate a workflow'
];

export default function OnboardingCompletePage() {
  return (
    <OnboardingShell activePath="/onboarding/complete">
      <div className="grid gap-6 xl:grid-cols-[1fr_0.85fr]">
        <section className="rounded-[2rem] border border-white/8 bg-gradient-to-br from-emerald-400/10 via-white/5 to-transparent p-8 shadow-glow backdrop-blur">
          <div className="inline-flex rounded-full border border-emerald-300/20 bg-emerald-300/10 px-3 py-1 text-xs uppercase tracking-[0.3em] text-emerald-100">
            Onboarding complete
          </div>
          <h2 className="mt-5 text-4xl font-semibold tracking-tight text-white md:text-6xl">
            Your first repo is connected.
          </h2>
          <p className="mt-5 max-w-2xl text-base leading-8 text-slate-300 md:text-lg">
            The platform has enough delivery evidence to build summaries, explain failures, and reason about changes across pipelines and deployments.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link href="/dashboard" className="rounded-full bg-cyan-300 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-200">
              Open dashboard
            </Link>
            <Link href="/ai-assistant" className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
              Ask a question
            </Link>
          </div>
        </section>

        <aside className="rounded-[2rem] border border-white/8 bg-white/5 p-6 backdrop-blur">
          <div className="text-xs uppercase tracking-[0.28em] text-slate-400">Next actions</div>
          <div className="mt-4 space-y-3">
            {nextActions.map((item) => (
              <div key={item} className="rounded-2xl border border-white/8 bg-black/15 p-4 text-sm text-slate-200">{item}</div>
            ))}
          </div>
        </aside>
      </div>
    </OnboardingShell>
  );
}

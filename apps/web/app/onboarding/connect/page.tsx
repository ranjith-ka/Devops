import Link from 'next/link';
import { OnboardingShell } from '@/components/onboarding-shell';

const providers = [
  { name: 'GitHub', description: 'Install the app or authorize via OAuth to access repos, PRs, Actions, and checks.' },
  { name: 'Azure DevOps', description: 'Connect projects, pipelines, repos, and deployment metadata with scoped permissions.' }
];

export default function RepositoryConnectPage() {
  return (
    <OnboardingShell activePath="/onboarding/connect">
      <div className="grid gap-6 xl:grid-cols-[1fr_0.9fr]">
        <section className="rounded-[2rem] border border-white/8 bg-white/5 p-6 shadow-glow backdrop-blur">
          <div className="text-xs uppercase tracking-[0.28em] text-cyan-200/80">Repository connection</div>
          <h2 className="mt-3 text-3xl font-semibold tracking-tight">Choose the provider and authorize the first repo sync.</h2>
          <div className="mt-6 grid gap-4 md:grid-cols-2">
            {providers.map((provider) => (
              <article key={provider.name} className="rounded-3xl border border-white/8 bg-black/20 p-5">
                <div className="text-xl font-semibold text-white">{provider.name}</div>
                <p className="mt-2 text-sm leading-7 text-slate-300">{provider.description}</p>
                <div className="mt-4 flex flex-wrap gap-2 text-xs text-slate-300">
                  <span className="rounded-full border border-white/10 px-3 py-1">Repo access</span>
                  <span className="rounded-full border border-white/10 px-3 py-1">Pipeline read</span>
                  <span className="rounded-full border border-white/10 px-3 py-1">PR metadata</span>
                </div>
              </article>
            ))}
          </div>
          <div className="mt-6 rounded-3xl border border-white/8 bg-black/20 p-5">
            <div className="text-sm font-medium text-slate-100">What will be imported</div>
            <div className="mt-3 grid gap-2 text-sm text-slate-300 md:grid-cols-2">
              {['commits', 'branches', 'pipelines', 'pipeline logs', 'PRs', 'deployment history', 'test results', 'security scans'].map((item) => (
                <div key={item} className="rounded-2xl border border-white/8 bg-white/5 px-3 py-2">{item}</div>
              ))}
            </div>
          </div>
          <div className="mt-6 flex flex-wrap gap-3">
            <Link href="/onboarding/sync" className="rounded-full bg-cyan-300 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-200">
              Authorize and sync
            </Link>
            <Link href="/onboarding/organization" className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
              Back
            </Link>
          </div>
        </section>

        <aside className="space-y-4">
          <div className="rounded-3xl border border-emerald-300/20 bg-emerald-400/10 p-5">
            <div className="text-sm font-medium text-emerald-100">Permissions are explicit</div>
            <p className="mt-2 text-sm leading-7 text-slate-300">
              The app only needs the scopes required to import repo and delivery data. Advanced automation stays locked until the user approves it.
            </p>
          </div>
          <div className="rounded-3xl border border-white/8 bg-white/5 p-5 text-sm leading-7 text-slate-300 backdrop-blur">
            Support single repo or multi repo sync, dry-run validation, and a visible resume state if the user leaves midway.
          </div>
        </aside>
      </div>
    </OnboardingShell>
  );
}

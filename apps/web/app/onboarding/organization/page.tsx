import Link from 'next/link';
import { OnboardingShell } from '@/components/onboarding-shell';

const plans = [
  { name: 'Trial', details: 'Fast start with one repo and limited sync depth.' },
  { name: 'Team', details: 'Multiple repos, shared insights, and collaboration.' },
  { name: 'Enterprise', details: 'RBAC, audit, advanced policy controls, and SSO.' }
];

export default function OrganizationSetupPage() {
  return (
    <OnboardingShell activePath="/onboarding/organization">
      <div className="grid gap-6 xl:grid-cols-[1fr_0.8fr]">
        <section className="rounded-[2rem] border border-white/8 bg-white/5 p-6 shadow-glow backdrop-blur">
          <div className="text-xs uppercase tracking-[0.28em] text-cyan-200/80">Organization setup</div>
          <h2 className="mt-3 text-3xl font-semibold tracking-tight">Create the workspace that owns your delivery data.</h2>
          <div className="mt-6 grid gap-4 md:grid-cols-2">
            <label className="space-y-2 text-sm text-slate-300">
              <span>Organization name</span>
              <input className="w-full rounded-2xl border border-white/10 bg-black/20 px-4 py-3 text-white outline-none ring-0 placeholder:text-slate-500" placeholder="Acme Engineering" />
            </label>
            <label className="space-y-2 text-sm text-slate-300">
              <span>Slug</span>
              <input className="w-full rounded-2xl border border-white/10 bg-black/20 px-4 py-3 text-white outline-none ring-0 placeholder:text-slate-500" placeholder="acme" />
            </label>
            <label className="space-y-2 text-sm text-slate-300">
              <span>Primary region</span>
              <select className="w-full rounded-2xl border border-white/10 bg-black/20 px-4 py-3 text-white outline-none ring-0">
                <option>us-east</option>
                <option>eu-west</option>
                <option>ap-south</option>
              </select>
            </label>
            <label className="space-y-2 text-sm text-slate-300">
              <span>Team size</span>
              <select className="w-full rounded-2xl border border-white/10 bg-black/20 px-4 py-3 text-white outline-none ring-0">
                <option>1-10</option>
                <option>11-50</option>
                <option>51-250</option>
                <option>250+</option>
              </select>
            </label>
          </div>
          <div className="mt-6 flex flex-wrap gap-3">
            <Link href="/onboarding/connect" className="rounded-full bg-cyan-300 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-200">
              Continue to repo connection
            </Link>
            <Link href="/onboarding" className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
              Back
            </Link>
          </div>
        </section>

        <aside className="space-y-4">
          {plans.map((plan) => (
            <article key={plan.name} className="rounded-3xl border border-white/8 bg-white/5 p-5 backdrop-blur">
              <div className="text-lg font-semibold">{plan.name}</div>
              <p className="mt-2 text-sm leading-7 text-slate-300">{plan.details}</p>
            </article>
          ))}
        </aside>
      </div>
    </OnboardingShell>
  );
}

import Link from 'next/link';
import { OnboardingShell } from '@/components/onboarding-shell';

const highlights = [
  'See value before you configure anything.',
  'Keep the flow linear and resumable.',
  'Show exactly what repo access is needed.',
  'Land on a populated dashboard after sync.'
];

export default function OnboardingWelcomePage() {
  return (
    <OnboardingShell activePath="/onboarding">
      <div className="grid gap-6 xl:grid-cols-[1.05fr_0.95fr]">
        <section className="rounded-[2rem] border border-white/8 bg-gradient-to-br from-cyan-400/12 via-white/5 to-transparent p-8 shadow-glow backdrop-blur">
          <div className="inline-flex rounded-full border border-cyan-300/20 bg-cyan-300/10 px-3 py-1 text-xs uppercase tracking-[0.3em] text-cyan-100">
            Welcome
          </div>
          <h2 className="mt-5 max-w-2xl text-4xl font-semibold tracking-tight text-white md:text-6xl">
            Connect one repository and turn delivery noise into evidence.
          </h2>
          <p className="mt-5 max-w-2xl text-base leading-8 text-slate-300 md:text-lg">
            Pipeline OS becomes useful the moment it can see your repository, pipelines, and deployments. This onboarding flow sets up the organization, connects the repo provider, and imports the first delivery signals.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link href="/onboarding/organization" className="rounded-full bg-cyan-300 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-200">
              Get started
            </Link>
            <Link href="/onboarding/connect" className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
              Connect repository
            </Link>
          </div>
        </section>

        <aside className="space-y-4">
          {highlights.map((item) => (
            <div key={item} className="rounded-3xl border border-white/8 bg-white/5 p-5 text-sm leading-7 text-slate-300 backdrop-blur">
              {item}
            </div>
          ))}
        </aside>
      </div>
    </OnboardingShell>
  );
}

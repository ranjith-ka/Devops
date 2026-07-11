import Link from 'next/link';
import type { ReactNode } from 'react';

const steps = [
  { href: '/onboarding', label: 'Welcome' },
  { href: '/onboarding/organization', label: 'Organization' },
  { href: '/onboarding/connect', label: 'Connect repo' },
  { href: '/onboarding/sync', label: 'Sync' },
  { href: '/onboarding/complete', label: 'Done' }
];

export function OnboardingShell({ children, activePath }: { children: ReactNode; activePath: string }) {
  return (
    <div className="min-h-screen text-slate-50">
      <div className="mx-auto min-h-screen max-w-6xl px-6 py-8">
        <div className="flex items-center justify-between gap-4">
          <div>
            <div className="text-xs uppercase tracking-[0.3em] text-cyan-300/80">Pipeline OS onboarding</div>
            <h1 className="mt-2 text-2xl font-semibold">Connect your delivery systems</h1>
          </div>
          <Link href="/" className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-100 transition hover:bg-white/10">
            Back to landing
          </Link>
        </div>

        <div className="mt-8 grid gap-4 rounded-[2rem] border border-white/8 bg-white/5 p-4 backdrop-blur md:grid-cols-5">
          {steps.map((step, index) => {
            const active = activePath === step.href;
            const done = steps.findIndex((item) => item.href === activePath) > index;
            return (
              <Link
                key={step.href}
                href={step.href as any}
                className={[
                  'rounded-2xl border px-4 py-3 text-sm transition',
                  active ? 'border-cyan-300/30 bg-cyan-400/15 text-white' : 'border-white/8 bg-black/15 text-slate-300 hover:bg-white/10',
                  done ? 'border-emerald-300/20 bg-emerald-400/10 text-emerald-100' : ''
                ].join(' ')}
              >
                <div className="text-xs uppercase tracking-[0.25em] opacity-70">Step {index + 1}</div>
                <div className="mt-1 font-medium">{step.label}</div>
              </Link>
            );
          })}
        </div>

        <div className="mt-8">{children}</div>
      </div>
    </div>
  );
}

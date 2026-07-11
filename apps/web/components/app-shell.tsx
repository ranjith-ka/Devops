import type { ReactNode } from 'react';
import Link from 'next/link';
import { navigation } from '@/lib/navigation';

export function AppShell({ children, activePath }: { children: ReactNode; activePath: string }) {
  return (
    <div className="min-h-screen text-slate-50">
      <div className="mx-auto flex min-h-screen max-w-[1600px]">
        <aside className="hidden w-72 shrink-0 border-r border-white/8 bg-white/3 px-5 py-6 backdrop-blur xl:block">
          <div className="mb-8">
            <div className="text-xs uppercase tracking-[0.28em] text-cyan-300/80">AI CI/CD Intelligence</div>
            <div className="mt-2 text-2xl font-semibold">Pipeline OS</div>
            <p className="mt-3 max-w-xs text-sm leading-6 text-slate-300">
              AI-native delivery intelligence for pipelines, deployments, logs, security, and cost.
            </p>
          </div>
          <nav className="space-y-1 text-sm text-slate-300">
            {navigation.map((item) => {
              const isActive = activePath === item.href;
              return (
                <Link
                  key={item.label}
                  href={item.href as any}
                  className={[
                    'block rounded-xl px-3 py-2 transition',
                    isActive ? 'bg-cyan-400/15 text-white ring-1 ring-cyan-300/20' : 'hover:bg-white/5 hover:text-white'
                  ].join(' ')}
                >
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </aside>
        <main className="flex-1">
          <div className="border-b border-white/8 bg-black/15 px-6 py-4 backdrop-blur">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <div className="text-xs uppercase tracking-[0.3em] text-slate-400">Enterprise SaaS</div>
                <h1 className="mt-1 text-xl font-semibold">Delivery intelligence and automation</h1>
              </div>
              <div className="flex items-center gap-2 text-sm text-slate-300">
                <span className="rounded-full border border-emerald-400/30 bg-emerald-400/10 px-3 py-1 text-emerald-300">Healthy</span>
                <span className="rounded-full border border-cyan-400/30 bg-cyan-400/10 px-3 py-1 text-cyan-200">AI online</span>
              </div>
            </div>
          </div>
          <div className="px-6 py-6">{children}</div>
        </main>
      </div>
    </div>
  );
}

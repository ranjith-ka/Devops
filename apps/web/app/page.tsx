import Link from 'next/link';
import { AppShell } from '@/components/app-shell';

const metrics = [
  { label: 'Pipeline health', value: '97.8%' },
  { label: 'Deploy risk', value: 'Low' },
  { label: 'Median build', value: '8m 41s' },
  { label: 'Flaky tests', value: '12' }
];

const features = [
  'Reason over pipelines, PRs, logs, deployments, metrics, and security findings.',
  'Generate pipeline, Terraform, Kubernetes, and Helm artifacts with policy guardrails.',
  'Predict failure probability, deployment duration, rollback risk, and cost drift.',
  'Support GPT, Claude, Llama, Qwen, and DeepSeek through a provider abstraction layer.'
];

const signals = [
  { label: 'Connected systems', value: 'GitHub, Azure DevOps, Kubernetes, Docker, Helm, Terraform' },
  { label: 'AI workflows', value: 'RAG, tool calling, memory, prompt templates, semantic search' },
  { label: 'Observability', value: 'OpenTelemetry, Prometheus, Grafana, Loki, Tempo' }
];

export default function HomePage() {
  return (
    <AppShell activePath="/">
      <div className="space-y-8">
        <section className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
          <div className="rounded-[2.25rem] border border-white/8 bg-gradient-to-br from-cyan-400/12 via-white/5 to-transparent p-8 shadow-glow backdrop-blur">
            <div className="inline-flex rounded-full border border-cyan-300/20 bg-cyan-300/10 px-3 py-1 text-xs uppercase tracking-[0.3em] text-cyan-100">
              AI-native CI/CD intelligence
            </div>
            <h2 className="mt-6 max-w-3xl text-4xl font-semibold tracking-tight text-white md:text-6xl">
              Understand every delivery decision before it becomes a production incident.
            </h2>
            <p className="mt-5 max-w-2xl text-base leading-8 text-slate-300 md:text-lg">
              Pipeline OS is a commercial-grade delivery intelligence platform that reasons over CI/CD, logs, metrics, PRs, security, and infrastructure so teams can troubleshoot, optimize, and automate with evidence instead of guesswork.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link href="/onboarding" className="rounded-full bg-cyan-300 px-5 py-3 text-sm font-semibold text-slate-950 transition hover:bg-cyan-200">
                Get started
              </Link>
              <Link href="/pipeline-generator" className="rounded-full border border-white/12 bg-white/5 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10">
                Generate workflow
              </Link>
            </div>
          </div>

          <div className="rounded-[2.25rem] border border-white/8 bg-white/5 p-6 shadow-glow backdrop-blur">
            <div className="text-xs uppercase tracking-[0.28em] text-slate-400">Live platform signals</div>
            <div className="mt-5 grid gap-4 sm:grid-cols-2">
              {metrics.map((metric) => (
                <article key={metric.label} className="rounded-3xl border border-white/8 bg-black/15 p-5">
                  <div className="text-sm text-slate-400">{metric.label}</div>
                  <div className="mt-2 text-2xl font-semibold tracking-tight text-white">{metric.value}</div>
                </article>
              ))}
            </div>
            <div className="mt-5 rounded-3xl border border-white/8 bg-black/20 p-5 text-sm leading-7 text-slate-300">
              Local-ready. Spec-first. Built for enterprise delivery workflows with explicit evidence, approval gates, and AI-generated artifacts.
            </div>
          </div>
        </section>

        <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {features.map((feature) => (
            <article key={feature} className="rounded-[1.75rem] border border-white/8 bg-white/5 p-5 text-sm leading-7 text-slate-300 backdrop-blur">
              {feature}
            </article>
          ))}
        </section>

        <section className="grid gap-6 xl:grid-cols-[0.95fr_1.05fr]">
          <article className="rounded-[2rem] border border-white/8 bg-white/5 p-6 backdrop-blur">
            <div className="text-xs uppercase tracking-[0.3em] text-slate-400">System coverage</div>
            <div className="mt-4 space-y-3 text-sm leading-7 text-slate-300">
              {signals.map((signal) => (
                <div key={signal.label} className="rounded-2xl border border-white/8 bg-black/15 p-4">
                  <div className="text-slate-100">{signal.label}</div>
                  <div className="mt-1 text-slate-300">{signal.value}</div>
                </div>
              ))}
            </div>
          </article>

          <article className="rounded-[2rem] border border-white/8 bg-white/5 p-6 backdrop-blur">
            <div className="text-xs uppercase tracking-[0.3em] text-slate-400">What the AI does</div>
            <div className="mt-4 grid gap-3 md:grid-cols-2">
              {['Root cause analysis', 'Pipeline explanation', 'Security review', 'Cost optimization', 'Release summarization', 'Incident summarization', 'Workflow generation', 'Deployment risk prediction'].map((item) => (
                <div key={item} className="rounded-2xl border border-white/8 bg-black/15 p-4 text-sm text-slate-200">
                  {item}
                </div>
              ))}
            </div>
          </article>
        </section>
      </div>
    </AppShell>
  );
}

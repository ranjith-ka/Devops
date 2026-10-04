'use client';
import { useState, type FormEvent } from 'react';
import Link from 'next/link';
import { ArrowRight, Box, Download, FileCode2, GitPullRequest, Layers, Package } from 'lucide-react';
import { defaultInput, generateFiles, validateInput, type OnboardingInput, type GeneratedFile } from '@/lib/flux-onboarding';
const fieldClass = 'mt-2 w-full rounded-xl border border-slate-700 bg-slate-950/70 px-3 py-2.5 text-sm text-white outline-none focus:border-cyan-400 focus:ring-2 focus:ring-cyan-400/20';
export function PlatformOnboarding() {
  const [input, setInput] = useState<OnboardingInput>(defaultInput);
  const [files, setFiles] = useState<GeneratedFile[]>([]);
  const [selected, setSelected] = useState(0);
  const [errors, setErrors] = useState<string[]>([]);
  const [downloading, setDownloading] = useState(false);
  const [downloaded, setDownloaded] = useState(false);
  function update<K extends keyof OnboardingInput>(key: K, value: OnboardingInput[K]) {
    setInput(previous => ({ ...previous, [key]: value })); setFiles([]); setErrors([]); setDownloaded(false);
  }
  function generate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); const issues = validateInput(input); setErrors(issues);
    if (issues.length) return;
    setFiles(generateFiles(input)); setSelected(0); setDownloaded(false);
  }
  async function download() {
    setDownloading(true); setErrors([]);
    try {
      const response = await fetch('/api/platform-onboarding', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(input) });
      if (!response.ok) { const result = await response.json(); throw new Error(result.error || 'Could not generate the download.'); }
      const url = URL.createObjectURL(await response.blob()); const anchor = document.createElement('a'); anchor.href = url;
      anchor.download = `${input.name}-${input.environment}-flux.zip`; anchor.click(); setTimeout(() => URL.revokeObjectURL(url), 1000); setDownloaded(true);
    } catch (error) { setErrors([error instanceof Error ? error.message : 'Download failed. Please try again.']); }
    finally { setDownloading(false); }
  }
  return <main className="mx-auto max-w-7xl px-5 pb-16 sm:px-8">
    <header className="flex flex-wrap items-center justify-between gap-4 border-b border-slate-700/60 py-6">
      <Link href="/" className="flex items-center gap-3"><span className="rounded-xl bg-cyan-400/15 p-2.5 text-cyan-300"><Layers size={22} aria-hidden="true" /></span><span className="font-semibold">Pipeline OS <span className="ml-2 font-normal text-slate-400">/ Flux platform</span></span></Link>
      <span className="rounded-full border border-slate-700 px-3 py-1 text-xs text-slate-300">Git-based onboarding</span>
    </header>
    <section className="py-10 sm:py-12"><div className="mb-4 flex items-center gap-2 text-sm text-cyan-300"><Box size={16} aria-hidden="true" /> Application onboarding</div>
      <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">Your next application, ready for Git.</h1>
      <p className="mt-4 max-w-2xl text-base leading-7 text-slate-300">Choose an image and environment. Review your deployment configuration, download it, and open a pull request. Flux takes over after merge.</p>
      <ol className="mt-7 flex flex-wrap gap-x-7 gap-y-3 text-sm text-slate-300">{['Configure application', 'Review files', 'Download & open PR'].map((step, index) => <li key={step} className="flex items-center gap-2"><span className="flex h-6 w-6 items-center justify-center rounded-full border border-slate-600 text-xs">{index + 1}</span>{step}</li>)}</ol>
    </section>
    <div className="grid items-start gap-6 lg:grid-cols-[0.9fr_1.1fr]">
      <form onSubmit={generate} noValidate className="rounded-2xl border border-slate-700/70 bg-slate-900/70 p-6 sm:p-7">
        <h2 className="flex items-center gap-2 text-lg font-semibold"><Package size={19} className="text-cyan-300" aria-hidden="true" /> Application details</h2><p className="mt-2 text-sm text-slate-400">Uses the repository’s shared application template.</p>
        <div className="mt-6 grid gap-5 sm:grid-cols-2">
          <label htmlFor="app-name" className="text-sm">Application name<input id="app-name" className={fieldClass} value={input.name} onChange={event => update('name', event.target.value)} autoComplete="off" placeholder="payments-api" maxLength={40} required /><span className="mt-1.5 block text-xs text-slate-400">Lowercase letters, numbers and hyphens.</span></label>
          <label htmlFor="environment" className="text-sm">Environment<select id="environment" className={fieldClass} value={input.environment} onChange={event => update('environment', event.target.value)}><option value="dev">Development</option><option value="staging">Staging</option><option value="production">Production</option></select></label>
          <label htmlFor="image-repository" className="text-sm sm:col-span-2">Image repository<input id="image-repository" className={fieldClass} value={input.repository} onChange={event => update('repository', event.target.value)} placeholder="ghcr.io/team/payments-api" autoComplete="off" required /><span className="mt-1.5 block text-xs text-slate-400">An existing image compatible with the shared chart’s serve command.</span></label>
          <label htmlFor="image-version" className="text-sm">Image version<input id="image-version" className={fieldClass} value={input.tag} onChange={event => update('tag', event.target.value)} placeholder="1.2.0" autoComplete="off" required /><span className="mt-1.5 block text-xs text-slate-400">Use a published versioned tag.</span></label>
          <label htmlFor="replicas" className="text-sm">Replicas<input id="replicas" className={fieldClass} type="number" min={1} max={20} step={1} value={Number.isNaN(input.replicas) ? '' : input.replicas} onChange={event => update('replicas', event.target.valueAsNumber)} required /></label>
          <label htmlFor="port" className="text-sm">Container port<input id="port" className={fieldClass} type="number" min={1} max={65535} step={1} value={Number.isNaN(input.port) ? '' : input.port} onChange={event => update('port', event.target.valueAsNumber)} required /></label>
          <label htmlFor="hostname" className="text-sm">Hostname <span className="text-slate-400">(optional)</span><input id="hostname" className={fieldClass} value={input.hostname} onChange={event => update('hostname', event.target.value)} placeholder="payments.example.com" autoComplete="off" /><span className="mt-1.5 block text-xs text-slate-400">Leave blank for internal access.</span></label>
        </div>
        {errors.length > 0 && <div role="alert" className="mt-5 rounded-xl border border-rose-400/30 bg-rose-400/10 p-4 text-sm text-rose-200"><ul className="list-inside list-disc space-y-2">{errors.map(error => <li key={error}>{error}</li>)}</ul></div>}
        <button type="submit" className="mt-7 flex w-full items-center justify-center gap-2 rounded-xl bg-cyan-300 px-4 py-3 text-sm font-semibold text-slate-950 hover:bg-cyan-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-cyan-300">Generate configuration <ArrowRight size={17} aria-hidden="true" /></button>
      </form>
      <section className="overflow-hidden rounded-2xl border border-slate-700/70 bg-slate-950/70" aria-label="Configuration preview">
        <div className="flex items-center justify-between gap-3 border-b border-slate-700/60 p-5"><h2 className="flex items-center gap-2 font-semibold"><FileCode2 size={18} className="text-cyan-300" aria-hidden="true" /> Review configuration</h2><span className="text-xs text-slate-400">{files.length ? `${files.length} files` : 'Awaiting configuration'}</span></div>
        {!files.length ? <div className="flex min-h-[400px] flex-col items-center justify-center px-8 text-center"><div className="rounded-2xl border border-slate-700 bg-slate-800/40 p-4 text-slate-400"><GitPullRequest size={30} aria-hidden="true" /></div><h3 className="mt-5 font-medium text-slate-200">A reviewable deployment starts here</h3><p className="mt-3 max-w-sm text-sm leading-6 text-slate-400">Generate your configuration to preview the values, namespace and Flux resources before downloading.</p></div> : <>
          <div className="flex flex-wrap gap-2 border-b border-slate-700/60 p-4" role="group" aria-label="Generated files">{files.map((file, index) => <button key={file.path} type="button" aria-pressed={selected === index} onClick={() => setSelected(index)} className={`rounded-lg px-3 py-2 text-xs focus-visible:outline focus-visible:outline-2 focus-visible:outline-cyan-300 ${selected === index ? 'bg-cyan-300/15 text-cyan-200' : 'text-slate-400 hover:bg-slate-800 hover:text-white'}`}>{file.path.split('/').pop()}{index === 0 ? ' (values)' : ''}</button>)}</div>
          <div className="break-all px-5 pt-4 text-xs leading-5 text-slate-400">{files[selected].path}</div><pre className="max-h-[460px] min-h-[280px] overflow-auto p-5 text-xs leading-6 text-slate-200"><code>{files[selected].content}</code></pre>
          <div className="border-t border-slate-700/60 p-5"><button type="button" disabled={downloading} onClick={download} className="flex w-full items-center justify-center gap-2 rounded-xl border border-cyan-300/40 bg-cyan-300/10 px-4 py-3 text-sm font-semibold text-cyan-200 hover:bg-cyan-300/20 disabled:cursor-wait disabled:opacity-60"><Download size={17} aria-hidden="true" />{downloading ? 'Preparing bundle…' : 'Download configuration ZIP'}</button><p className="mt-3 text-center text-xs text-slate-400" role="status">{downloaded ? 'Download ready. Follow ONBOARDING.md to prepare your PR.' : 'Includes repository paths and PR instructions.'}</p></div>
        </>}
      </section>
    </div>
    <section className="mt-8 grid gap-4 sm:grid-cols-3" aria-label="Deployment workflow">{[['Review in Git', 'Extract the bundle, add its resource entry, and open a PR.'], ['Merge to deploy', 'Flux applies the configuration and Helm installs the application.'], ['Upgrade by version', 'Change the image tag in your values file and merge another PR.']].map(([title, description]) => <div key={title} className="border-t border-slate-700/60 pt-5"><h3 className="text-sm font-medium">{title}</h3><p className="mt-2 text-sm leading-6 text-slate-400">{description}</p></div>)}</section>
  </main>;
}

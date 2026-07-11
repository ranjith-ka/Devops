export function ScreenSpec({
  title,
  subtitle,
  api,
  wireframe,
  hierarchy,
  journey,
  stateModel,
  states,
  evidence,
}: {
  title: string;
  subtitle: string;
  api: string[];
  wireframe: string;
  hierarchy: string[];
  journey: string[];
  stateModel: string[];
  states: { label: string; description: string }[];
  evidence: string[];
}) {
  return (
    <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
      <section className="space-y-6">
        <article className="rounded-[2rem] border border-white/8 bg-white/5 p-6 shadow-glow backdrop-blur">
          <div className="text-xs uppercase tracking-[0.28em] text-cyan-200/80">Screen</div>
          <h2 className="mt-2 text-3xl font-semibold tracking-tight">{title}</h2>
          <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-300">{subtitle}</p>
        </article>
        <article className="rounded-[2rem] border border-white/8 bg-black/20 p-6 backdrop-blur">
          <div className="text-sm font-medium text-slate-100">Wireframe</div>
          <pre className="mt-4 overflow-x-auto rounded-2xl border border-white/8 bg-slate-950/80 p-4 text-xs leading-6 text-slate-300">{wireframe}</pre>
        </article>
        <div className="grid gap-4 md:grid-cols-2">
          <article className="rounded-3xl border border-white/8 bg-white/5 p-5 backdrop-blur">
            <div className="text-sm font-medium text-slate-100">Component Hierarchy</div>
            <ul className="mt-3 space-y-2 text-sm leading-6 text-slate-300">
              {hierarchy.map((item) => <li key={item}>{item}</li>)}
            </ul>
          </article>
          <article className="rounded-3xl border border-white/8 bg-white/5 p-5 backdrop-blur">
            <div className="text-sm font-medium text-slate-100">User Journey</div>
            <ul className="mt-3 space-y-2 text-sm leading-6 text-slate-300">
              {journey.map((item) => <li key={item}>{item}</li>)}
            </ul>
          </article>
        </div>
      </section>
      <aside className="space-y-6">
        <article className="rounded-[2rem] border border-white/8 bg-white/5 p-6 backdrop-blur">
          <div className="text-sm font-medium text-slate-100">API Mapping</div>
          <div className="mt-4 flex flex-wrap gap-2">
            {api.map((item) => <span key={item} className="rounded-full border border-cyan-400/20 bg-cyan-400/10 px-3 py-1 text-xs text-cyan-100">{item}</span>)}
          </div>
        </article>
        <article className="rounded-[2rem] border border-white/8 bg-white/5 p-6 backdrop-blur">
          <div className="text-sm font-medium text-slate-100">State Model</div>
          <ul className="mt-3 space-y-2 text-sm leading-6 text-slate-300">
            {stateModel.map((item) => <li key={item}>{item}</li>)}
          </ul>
        </article>
        <article className="rounded-[2rem] border border-white/8 bg-white/5 p-6 backdrop-blur">
          <div className="text-sm font-medium text-slate-100">Loading, Error, Empty</div>
          <div className="mt-4 space-y-3 text-sm text-slate-300">
            {states.map((state) => (
              <div key={state.label} className="rounded-2xl border border-white/8 bg-black/20 p-4">
                <div className="font-medium text-slate-100">{state.label}</div>
                <p className="mt-1 leading-6 text-slate-300">{state.description}</p>
              </div>
            ))}
          </div>
        </article>
        <article className="rounded-[2rem] border border-white/8 bg-gradient-to-br from-cyan-400/10 via-white/5 to-transparent p-6 backdrop-blur">
          <div className="text-sm font-medium text-slate-100">Evidence Sources</div>
          <ul className="mt-3 space-y-2 text-sm leading-6 text-slate-300">
            {evidence.map((item) => <li key={item}>{item}</li>)}
          </ul>
        </article>
      </aside>
    </div>
  );
}

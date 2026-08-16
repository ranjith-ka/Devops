const repo = "https://github.com/ranjith-ka/Devops";

const services = [
  {
    number: "01",
    title: "Platform engineering",
    text: "Design paved roads for developers with Kubernetes, GitOps, reusable infrastructure and secure runtime standards.",
    tags: ["Kubernetes", "FluxCD", "Terraform"],
  },
  {
    number: "02",
    title: "Delivery modernization",
    text: "Move fragmented CI/CD into dependable delivery systems with faster feedback, clear guardrails and measurable outcomes.",
    tags: ["GitHub Actions", "Go", "Developer Experience"],
  },
  {
    number: "03",
    title: "AI-ready operations",
    text: "Connect logs, metrics and traces to AI-assisted investigation, evaluation and production model infrastructure.",
    tags: ["Grafana", "LangGraph", "MCP"],
  },
];

const projects = [
  {
    title: "AI CI/CD Intelligence Platform",
    label: "Product blueprint",
    text: "A spec-first platform for pipeline intelligence, inference services and engineering decision support.",
    href: `${repo}/tree/main/product/ai-cicd-platform`,
  },
  {
    title: "BuildPilot",
    label: "Go · Kubernetes · AI",
    text: "A control plane and cluster agent exploring AI-assisted build analysis and platform automation.",
    href: `${repo}/tree/main/projects/buildpilot`,
  },
  {
    title: "RAG Documentation Assistant",
    label: "Retrieval architecture",
    text: "A practical documentation retrieval assistant covering ingestion, embeddings and grounded responses.",
    href: `${repo}/tree/main/tutorials/rag-documentation-assistant`,
  },
];

const tutorials = [
  {
    index: "001",
    title: "Kubernetes in Docker with KIND",
    topic: "Local platform engineering",
    href: `${repo}/tree/main/kind`,
  },
  {
    index: "002",
    title: "SPIFFE/SPIRE on Kubernetes",
    topic: "Workload identity",
    href: `${repo}/blob/main/kind/SPIFFE_SPIRE_TUTORIAL.md`,
  },
  {
    index: "003",
    title: ".NET for DevOps engineers",
    topic: "Application fundamentals",
    href: `${repo}/blob/main/tutorials/dotnet/DOTNET_FOR_DEVOPS.md`,
  },
  {
    index: "004",
    title: "PostgreSQL learning path",
    topic: "Data foundations",
    href: `${repo}/tree/main/tutorials/postgreSQL`,
  },
  {
    index: "005",
    title: "Go bootcamp exercises",
    topic: "Platform tooling",
    href: `${repo}/tree/main/tutorials/golangbootcamp`,
  },
  {
    index: "006",
    title: "Python practical notes",
    topic: "Automation basics",
    href: `${repo}/tree/main/tutorials/python`,
  },
];

export default function Home() {
  return (
    <main>
      <nav className="nav shell" aria-label="Main navigation">
        <a className="brand" href="#top" aria-label="Ranjith K A, home">
          <span className="brand-mark">RK</span>
          <span>Ranjith K A</span>
        </a>
        <div className="nav-links">
          <a href="#work">Work</a>
          <a href="#tutorials">Tutorials</a>
          <a href="#services">Services</a>
        </div>
        <a className="nav-cta" href="https://www.linkedin.com/in/ranjith-k-a-05980522/" target="_blank" rel="noreferrer">
          Let&apos;s talk <span aria-hidden="true">↗</span>
        </a>
      </nav>

      <section className="hero shell" id="top">
        <div className="hero-copy">
          <p className="eyebrow"><span /> Bengaluru, India · Open to collaboration</p>
          <h1>I build platforms that help engineering teams <em>move with confidence.</em></h1>
          <p className="hero-text">
            Senior Platform &amp; DevOps engineer with 15+ years of experience across Kubernetes,
            developer enablement, observability and production-ready AI infrastructure.
          </p>
          <div className="hero-actions">
            <a className="button primary" href="#work">Explore my work <span>↓</span></a>
            <a className="button ghost" href={repo} target="_blank" rel="noreferrer">View GitHub <span>↗</span></a>
          </div>
        </div>
        <aside className="status-card" aria-label="Current focus">
          <div className="card-top"><span>Current focus</span><span className="pulse" /></div>
          <p>AI-native platform engineering</p>
          <div className="terminal">
            <span><b>$</b> platform status</span>
            <span><i>✓</i> Kubernetes foundations</span>
            <span><i>✓</i> GitOps delivery</span>
            <span><i>✓</i> AI observability</span>
            <span><i>→</i> Agentic operations</span>
          </div>
          <div className="availability"><span>Available for</span><strong>Architecture · Advisory · Build</strong></div>
        </aside>
      </section>

      <div className="signal-strip">
        <div className="shell signal-inner">
          <span>PLATFORM ENGINEERING</span><b>◆</b><span>KUBERNETES</span><b>◆</b><span>AI INFRASTRUCTURE</span><b>◆</b><span>OBSERVABILITY</span><b>◆</b><span>DEVELOPER EXPERIENCE</span>
        </div>
      </div>

      <section className="section shell" id="services">
        <div className="section-heading">
          <div><p className="kicker">How I can help</p><h2>From infrastructure to<br/><em>engineering leverage.</em></h2></div>
          <p>I help teams turn complex infrastructure into an understandable, secure and productive internal platform.</p>
        </div>
        <div className="service-grid">
          {services.map((service) => (
            <article className="service-card" key={service.number}>
              <span className="service-number">{service.number}</span>
              <h3>{service.title}</h3>
              <p>{service.text}</p>
              <div className="tags">{service.tags.map((tag) => <span key={tag}>{tag}</span>)}</div>
            </article>
          ))}
        </div>
      </section>

      <section className="section project-section" id="work">
        <div className="shell">
          <div className="section-heading light">
            <div><p className="kicker">Selected work</p><h2>Ideas made <em>operational.</em></h2></div>
            <a href={repo} target="_blank" rel="noreferrer">All repositories <span>↗</span></a>
          </div>
          <div className="project-grid">
            {projects.map((project, index) => (
              <a className={`project-card project-${index + 1}`} href={project.href} target="_blank" rel="noreferrer" key={project.title}>
                <span className="project-label">{project.label}</span>
                <div><h3>{project.title}</h3><p>{project.text}</p></div>
                <span className="project-arrow">↗</span>
              </a>
            ))}
          </div>
        </div>
      </section>

      <section className="section shell" id="tutorials">
        <div className="section-heading tutorial-heading">
          <div><p className="kicker">Learn in public</p><h2>Field notes &amp; <em>tutorials.</em></h2></div>
          <p>Practical notes built while learning, debugging and shipping—kept open for the next engineer.</p>
        </div>
        <div className="tutorial-list">
          {tutorials.map((tutorial) => (
            <a href={tutorial.href} target="_blank" rel="noreferrer" className="tutorial-row" key={tutorial.index}>
              <span className="tutorial-index">{tutorial.index}</span>
              <h3>{tutorial.title}</h3>
              <span className="tutorial-topic">{tutorial.topic}</span>
              <span className="row-arrow">↗</span>
            </a>
          ))}
        </div>
      </section>

      <section className="cta-section">
        <div className="shell cta-inner">
          <p className="kicker">Have a platform challenge?</p>
          <h2>Let&apos;s turn it into a<br/><em>system your team trusts.</em></h2>
          <p>Available for architecture reviews, platform strategy, DevOps modernization and AI infrastructure advisory.</p>
          <div className="hero-actions centered">
            <a className="button primary" href="https://www.linkedin.com/in/ranjith-k-a-05980522/" target="_blank" rel="noreferrer">Start a conversation <span>↗</span></a>
            <a className="button ghost light-ghost" href={repo} target="_blank" rel="noreferrer">Follow on GitHub</a>
          </div>
        </div>
      </section>

      <footer className="footer shell">
        <div className="brand"><span className="brand-mark">RK</span><span>Ranjith K A</span></div>
        <p>Building dependable platforms for people who build software.</p>
        <div><a href={repo} target="_blank" rel="noreferrer">GitHub</a><a href="https://www.linkedin.com/in/ranjith-k-a-05980522/" target="_blank" rel="noreferrer">LinkedIn</a></div>
      </footer>
    </main>
  );
}

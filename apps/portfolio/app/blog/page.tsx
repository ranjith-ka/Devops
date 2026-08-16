import type { Metadata } from "next";
import Link from "next/link";
import { posts } from "./posts";

export const metadata: Metadata = {
  title: "Kubernetes Field Notes | Ranjith K A",
  description: "Hands-on KIND and Minikube tutorials covering local Kubernetes, Gateway API, workload identity, Flux and KEDA.",
};

export default function BlogPage() {
  return (
    <main className="blog-page">
      <nav className="nav shell" aria-label="Blog navigation">
        <Link className="brand" href="/"><span className="brand-mark">RK</span><span>Ranjith K A</span></Link>
        <div className="nav-links"><Link href="/">Profile</Link><a href="#articles">Articles</a><a href="https://github.com/ranjith-ka/Devops" target="_blank" rel="noreferrer">GitHub</a></div>
        <Link className="nav-cta" href="/">Back home <span>↗</span></Link>
      </nav>
      <header className="blog-hero shell">
        <p className="eyebrow"><span /> Field notes from real platform work</p>
        <h1>Kubernetes learning,<br/><em>made practical.</em></h1>
        <p>Hands-on tutorials for building, securing and operating local Kubernetes platforms with KIND and Minikube.</p>
        <div className="blog-count"><strong>{posts.length}</strong><span>published<br/>guides</span></div>
      </header>
      <section className="blog-index shell" id="articles">
        <div className="filter-row"><span>All articles</span><span>KIND · Minikube · GitOps · Security</span></div>
        <div className="article-grid">
          {posts.map((post, index) => (
            <Link className={`article-card accent-${post.accent}`} href={`/blog/${post.slug}`} key={post.slug}>
              <div className="article-meta"><span>{post.category}</span><span>{post.readTime}</span></div>
              <div className="article-visual"><span>{String(index + 1).padStart(2, "0")}</span><i>{post.category === "KIND" ? "K8s" : "mini"}</i></div>
              <div className="article-copy"><p>{post.level}</p><h2>{post.title}</h2><p>{post.excerpt}</p></div>
              <span className="article-link">Read guide <b>→</b></span>
            </Link>
          ))}
        </div>
      </section>
      <footer className="footer shell"><div className="brand"><span className="brand-mark">RK</span><span>Ranjith K A</span></div><p>Practical platform engineering, shared openly.</p><div><Link href="/">Profile</Link><a href="https://github.com/ranjith-ka/Devops" target="_blank" rel="noreferrer">GitHub</a></div></footer>
    </main>
  );
}

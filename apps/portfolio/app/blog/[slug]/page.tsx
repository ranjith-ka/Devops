import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { getPost, posts } from "../posts";

type Props = { params: Promise<{ slug: string }> };

export function generateStaticParams() { return posts.map((post) => ({ slug: post.slug })) }

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const post = getPost((await params).slug);
  return post ? { title: `${post.title} | Ranjith K A`, description: post.excerpt } : {};
}

export default async function ArticlePage({ params }: Props) {
  const post = getPost((await params).slug);
  if (!post) notFound();
  return (
    <main className="article-page">
      <nav className="nav shell" aria-label="Article navigation">
        <Link className="brand" href="/"><span className="brand-mark">RK</span><span>Ranjith K A</span></Link>
        <div className="nav-links"><Link href="/">Profile</Link><Link href="/blog">All articles</Link></div>
        <Link className="nav-cta" href="/blog">Blog index <span>↗</span></Link>
      </nav>
      <header className={`article-hero accent-${post.accent}`}>
        <div className="shell article-hero-inner">
          <div><Link className="back-link" href="/blog">← All field notes</Link><div className="article-meta"><span>{post.category}</span><span>{post.level}</span><span>{post.readTime}</span></div><h1>{post.title}</h1><p>{post.excerpt}</p></div>
          <div className="article-glyph"><span>{post.category === "KIND" ? "K8s" : "mini"}</span><small>{post.date}</small></div>
        </div>
      </header>
      <article className="article-body shell">
        <aside className="article-aside"><span>In this guide</span>{post.sections.map((section, index) => <a href={`#section-${index + 1}`} key={section.heading}>{String(index + 1).padStart(2, "0")} {section.heading}</a>)}</aside>
        <div className="article-content">
          {post.sections.map((section, index) => (
            <section id={`section-${index + 1}`} key={section.heading}><span className="section-no">{String(index + 1).padStart(2, "0")}</span><h2>{section.heading}</h2><p>{section.body}</p>{section.code && <pre><code>{section.code}</code></pre>}{section.note && <div className="article-note"><strong>Platform note</strong><p>{section.note}</p></div>}</section>
          ))}
          <div className="source-card"><div><span>Original field notes</span><p>View the source tutorial, manifests and configuration in the Devops repository.</p></div><a href={post.source} target="_blank" rel="noreferrer">Open on GitHub ↗</a></div>
          <div className="article-end"><p>Written from hands-on platform engineering experiments by</p><strong>Ranjith K A</strong><span>Platform Engineering · DevOps · AI Infrastructure</span></div>
        </div>
      </article>
      <footer className="footer shell"><div className="brand"><span className="brand-mark">RK</span><span>Ranjith K A</span></div><p>Building dependable platforms for people who build software.</p><div><Link href="/blog">More articles</Link><Link href="/">Profile</Link></div></footer>
    </main>
  );
}

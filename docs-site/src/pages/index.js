import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import styles from './index.module.css';

function Hero() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <h1 className="hero__title">{siteConfig.title}</h1>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <p className={styles.heroDescription}>
          One CLI to scaffold reproducible, containerized dev environments with
          built-in AI context, curated skills, and structured work processes.
        </p>
        <div className={styles.buttons}>
          <Link className="button button--secondary button--lg" to="/docs/getting-started/installation">
            Get Started →
          </Link>
        </div>
      </div>
    </header>
  );
}

const features = [
  {
    title: 'One Command Setup',
    description: 'From zero to a fully configured dev environment with aibox init. Container, AI context, skills, and theming — all scaffolded.',
  },
  {
    title: '5 Process Packages',
    description: 'minimal, managed, software, research, product — pick the one that fits your project and get the right set of agent skills installed automatically.',
  },
  {
    title: 'Curated Skills via processkit',
    description: 'Vetted AI agent skills following the open SKILL.md standard — handcrafted with examples and reference files, managed by processkit and installed by aibox.',
  },
  {
    title: '25 Add-ons',
    description: 'Python, Rust, Node, Go, LaTeX, Kubernetes, cloud providers, docs frameworks, and more — each with per-tool version selection from curated lists.',
  },
];

function Feature({title, description}) {
  return (
    <div className={clsx('col col--6')}>
      <div className="padding-horiz--md padding-vert--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

function Features() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {features.map((props, idx) => <Feature key={idx} {...props} />)}
        </div>
      </div>
    </section>
  );
}

function QuickStart() {
  return (
    <section className={styles.quickstart}>
      <div className="container">
        <h2>Quick Start</h2>
        <pre><code>{`curl -fsSL https://raw.githubusercontent.com/projectious-work/aibox/main/scripts/install.sh | bash
aibox init --name my-project
aibox sync && aibox start`}</code></pre>
      </div>
    </section>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout title="Home" description={siteConfig.tagline}>
      <Hero />
      <main>
        <Features />
        <QuickStart />
      </main>
    </Layout>
  );
}

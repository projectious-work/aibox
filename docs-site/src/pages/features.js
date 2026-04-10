import React from 'react';
import Layout from '@theme/Layout';

export default function Features() {
  return (
    <Layout title="Features" description="aibox features overview">
      <div className="container margin-vert--lg">
        <h1>Features</h1>

        <h2>Single Base Image + Add-ons</h2>
        <p>One base-debian image provides the foundation. Everything else — Python, Rust, Node, LaTeX, Kubernetes tools — is a declarative add-on with per-tool version selection.</p>

        <h2>5 Process Packages</h2>
        <p>minimal, managed, software, research, product — composable tiers that install the right skill set for your project type. Select in aibox.toml; switch any time with aibox sync.</p>

        <h2>Curated AI Skills via processkit</h2>
        <p>Skills follow the open SKILL.md standard with progressive disclosure: trigger conditions, structured instructions, and real examples. The skill catalogue is owned and versioned by <a href="https://github.com/projectious-work/processkit">processkit</a>; aibox handles installation and sync.</p>

        <h2>7 Color Themes</h2>
        <p>Gruvbox Dark, Catppuccin Mocha, Catppuccin Latte, Dracula, Tokyo Night, Nord, and Projectious — applied consistently across Zellij, Vim, Yazi, lazygit, and Starship.</p>

        <h2>Declarative Configuration</h2>
        <p>One aibox.toml controls everything: base image, add-ons with versions, process packages, skill selection, themes, and AI providers. Run aibox sync to reconcile.</p>

        <h2>Migration System</h2>
        <p>When aibox updates, migration documents are auto-generated with safety headers, action items, and verification checklists. Your AI agent picks them up at session start.</p>
      </div>
    </Layout>
  );
}

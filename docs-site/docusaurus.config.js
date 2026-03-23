import {themes as prismThemes} from 'prism-react-renderer';

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'aibox',
  tagline: 'AI-ready development environments, containerized',
  favicon: 'img/favicon.ico',
  url: 'https://projectious-work.github.io',
  baseUrl: '/aibox/',
  organizationName: 'projectious-work',
  projectName: 'aibox',
  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',
  i18n: { defaultLocale: 'en', locales: ['en'] },
  markdown: {
    format: 'detect',
  },
  presets: [['classic', {
    docs: {
      sidebarPath: './sidebars.js',
      editUrl: 'https://github.com/projectious-work/aibox/tree/main/docs-site/',
    },
    blog: false,
    theme: { customCss: './src/css/custom.css' },
  }]],
  themeConfig: {
    colorMode: { defaultMode: 'dark', respectPrefersColorScheme: true },
    navbar: {
      title: 'aibox',
      items: [
        { type: 'docSidebar', sidebarId: 'docs', position: 'left', label: 'Docs' },
        { to: '/features', label: 'Features', position: 'left' },
        { to: '/changelog', label: 'Changelog', position: 'left' },
        { href: 'https://github.com/projectious-work/aibox', label: 'GitHub', position: 'right' },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        { title: 'Docs', items: [
          { label: 'Getting Started', to: '/docs/getting-started/installation' },
          { label: 'Configuration', to: '/docs/reference/configuration' },
          { label: 'Skills Library', to: '/docs/skills/' },
        ]},
        { title: 'Project', items: [
          { label: 'GitHub', href: 'https://github.com/projectious-work/aibox' },
          { label: 'Changelog', to: '/changelog' },
        ]},
      ],
      copyright: `\u00a9 ${new Date().getFullYear()} projectious.work`,
    },
    prism: { theme: prismThemes.github, darkTheme: prismThemes.dracula },
  },
};

export default config;

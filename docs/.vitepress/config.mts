import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'promptext',
  description: 'Smart code context extractor',
  base: '/promptext/',
  appearance: false,
  cleanUrls: true,
  head: [
    ['link', { rel: 'preconnect', href: 'https://fonts.googleapis.com' }],
    ['link', { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' }],
    ['link', { href: 'https://fonts.googleapis.com/css2?family=Work+Sans:wght@300;400;700&family=JetBrains+Mono:wght@400;500&display=swap', rel: 'stylesheet' }],
  ],
  themeConfig: {
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Config', link: '/reference/configuration' },
      { text: 'Library', link: '/reference/library-usage' },
      { text: 'Changelog', link: '/changelog' },
    ],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/guide/getting-started' },
          { text: 'Output Formats', link: '/guide/output-formats' },
          { text: 'File Filtering', link: '/guide/file-filtering' },
          { text: 'Relevance Filtering', link: '/guide/relevance-filtering' },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'Configuration', link: '/reference/configuration' },
          { text: 'Library Usage', link: '/reference/library-usage' },
          { text: 'Performance', link: '/reference/performance' },
        ],
      },
      {
        text: 'Changelog',
        link: '/changelog',
      },
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/1broseidon/promptext' },
    ],
    outline: { level: [2, 3] },
  },
})

// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://promptext.sh',
	base: '/',
	integrations: [
		starlight({
			title: 'promptext',
			description: 'Smart code context extractor for AI assistants',
			logo: {
				light: './src/assets/logo-light.svg',
				dark: './src/assets/logo.svg',
			},
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/1broseidon/promptext' }
			],
			customCss: ['./src/styles/custom.css'],
			editLink: {
				baseUrl: 'https://github.com/1broseidon/promptext/edit/main/docs-astro/',
			},
			// Disable built-in homepage to use custom landing page
			disable404Route: false,
			sidebar: [
				{ label: 'Getting Started', slug: 'getting-started' },
				{ label: 'Configuration', slug: 'configuration' },
				{ label: 'File Filtering', slug: 'file-filtering' },
				{ label: 'Token Analysis', slug: 'token-analysis' },
				{ label: 'Project Analysis', slug: 'project-analysis' },
				{ label: 'Output Formats', slug: 'output-formats' },
				{ label: 'Performance', slug: 'performance' },
				{ label: 'Changelog', slug: 'changelog' },
			],
		}),
	],
});

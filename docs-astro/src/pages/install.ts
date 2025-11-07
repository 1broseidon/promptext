import type { APIRoute } from 'astro';
import { readFileSync } from 'fs';
import { join } from 'path';

export const GET: APIRoute = () => {
  try {
    const scriptPath = join(process.cwd(), 'public', 'scripts', 'install.sh');
    const script = readFileSync(scriptPath, 'utf-8');

    return new Response(script, {
      status: 200,
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
        'Cache-Control': 'public, max-age=3600',
        'X-Content-Type-Options': 'nosniff',
      },
    });
  } catch (error) {
    return new Response('Install script not found', {
      status: 404,
      headers: { 'Content-Type': 'text/plain' },
    });
  }
};

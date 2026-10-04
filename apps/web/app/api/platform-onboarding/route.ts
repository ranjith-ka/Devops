import { generateFiles, validateInput, type OnboardingInput } from '@/lib/flux-onboarding';
import { createZip } from '@/lib/zip';

export const runtime = 'nodejs';
export async function POST(request: Request) {
  if (Number(request.headers.get('content-length')) > 8192) return Response.json({ error: 'Request is too large.' }, { status: 413 });
  let input: OnboardingInput;
  try {
    const body: unknown = await request.json();
    if (!body || typeof body !== 'object' || Array.isArray(body)) throw new Error();
    const value = body as Record<string, unknown>;
    for (const key of ['name', 'environment', 'repository', 'tag', 'hostname']) if (typeof value[key] !== 'string') throw new Error();
    if (typeof value.replicas !== 'number' || typeof value.port !== 'number') throw new Error();
    input = value as OnboardingInput;
  } catch { return Response.json({ error: 'Provide a valid onboarding form.' }, { status: 400 }); }
  const errors = validateInput(input);
  if (errors.length) return Response.json({ error: errors.join(' ') }, { status: 400 });
  return new Response(Uint8Array.from(createZip(generateFiles(input))), {
    headers: { 'Content-Type': 'application/zip', 'Content-Disposition': `attachment; filename="${input.name}-${input.environment}-flux.zip"`, 'Cache-Control': 'no-store' },
  });
}

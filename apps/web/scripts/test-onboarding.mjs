import { readFileSync } from 'node:fs';
import assert from 'node:assert/strict';
import test from 'node:test';
import ts from 'typescript';
async function load(path) {
  const source = readFileSync(new URL(path, import.meta.url), 'utf8');
  const { outputText } = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2022 } });
  return import(`data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`);
}
const { defaultInput, generateFiles, validateInput } = await load('../lib/flux-onboarding.ts');
const { createZip } = await load('../lib/zip.ts');
test('reject invalid names, mutable tags and invalid deployment inputs', () => {
  for (const change of [{ name: '../escape' }, { name: 'MixedCase' }, { tag: 'latest' }, { repository: 'repo:version' }, { port: 0 }, { replicas: 0 }, { environment: 'unknown' }, { hostname: 'https://example.com' }]) {
    assert.ok(validateInput({ ...defaultInput, ...change }).length, JSON.stringify(change));
    assert.throws(() => generateFiles({ ...defaultInput, ...change }));
  }
});
test('bundle isolates each environment and references its values and artifact', () => {
  const files = generateFiles({ ...defaultInput, name: 'payments', environment: 'production', hostname: 'payments.example.com' });
  assert.equal(files.length, 6);
  assert.equal(new Set(files.map(file => file.path)).size, 6);
  const values = files[0]; assert.match(values.content, /tag: "0.0.1"/); assert.match(values.content, /host: "payments.example.com"/);
  const generator = files.find(file => file.path.endsWith('artifact-generator.yaml'));
  assert.ok(generator.content.includes(`@repo/${values.path}`));
  assert.match(generator.content, /namespace: "default"/);
  const release = files.find(file => file.path.endsWith('helm-release.yaml'));
  assert.match(release.content, /namespace: "payments-production"/); assert.match(release.content, /kind: "ExternalArtifact"/);
  assert.ok(files.at(-1).content.includes('../onboarded/payments-production'));
});
test('changing version changes values without introducing a HelmRelease override', () => {
  const before = generateFiles(defaultInput), after = generateFiles({ ...defaultInput, tag: '1.2.0' });
  assert.notEqual(before[0].content, after[0].content);
  assert.equal(before[3].content, after[3].content);
  assert.ok(!after[3].content.includes('image:'));
});
test('ZIP stores all generated files with exact contents and central directory entries', () => {
  const files = generateFiles(defaultInput), zip = createZip(files);
  let offset = 0;
  for (const file of files) {
    assert.equal(zip.readUInt32LE(offset), 0x04034b50);
    const size = zip.readUInt32LE(offset + 18), length = zip.readUInt16LE(offset + 26);
    assert.equal(zip.subarray(offset + 30, offset + 30 + length).toString(), file.path);
    assert.equal(zip.subarray(offset + 30 + length, offset + 30 + length + size).toString(), file.content);
    offset += 30 + length + size;
  }
  assert.equal(zip.readUInt32LE(offset), 0x02014b50);
  assert.equal(zip.readUInt32LE(zip.length - 22), 0x06054b50);
  assert.equal(zip.readUInt16LE(zip.length - 12), files.length);
});

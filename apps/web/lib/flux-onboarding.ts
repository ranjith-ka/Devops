export type OnboardingInput = {
  name: string; environment: string; repository: string; tag: string;
  replicas: number; port: number; hostname: string;
};
export type GeneratedFile = { path: string; content: string };
export const defaultInput: OnboardingInput = {
  name: 'hello-world', environment: 'dev', repository: 'ranjithka/canary',
  tag: '0.0.1', replicas: 1, port: 8080, hostname: '',
};
export function validateInput(input: OnboardingInput): string[] {
  const errors: string[] = [];
  if (!/^[a-z][a-z0-9-]{0,38}[a-z0-9]$|^[a-z]$/.test(input.name)) errors.push('Use an application name of 1–40 lowercase letters, numbers or hyphens, starting with a letter and ending with a letter or number.');
  if (!['dev', 'staging', 'production'].includes(input.environment)) errors.push('Choose dev, staging or production.');
  if (!/^(?:[a-z0-9.-]+(?::[0-9]+)?\/)?[a-z0-9]+(?:[._-][a-z0-9]+)*(?:\/[a-z0-9]+(?:[._-][a-z0-9]+)*)*$/.test(input.repository)) errors.push('Enter an image repository without a tag or URL scheme.');
  if (!/^[a-zA-Z0-9_][a-zA-Z0-9_.-]{0,127}$/.test(input.tag) || input.tag === 'latest') errors.push('Use a versioned image tag, such as 0.0.1, rather than latest.');
  if (!Number.isInteger(input.replicas) || input.replicas < 1 || input.replicas > 20) errors.push('Replicas must be between 1 and 20.');
  if (!Number.isInteger(input.port) || input.port < 1 || input.port > 65535) errors.push('Container port must be between 1 and 65535.');
  if (input.hostname && (input.hostname.length > 253 || !input.hostname.split('.').every(part => /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/.test(part)))) errors.push('Enter a valid lowercase hostname without a scheme or path.');
  return errors;
}
function yaml(value: unknown, indent = 0): string {
  const pad = ' '.repeat(indent);
  if (Array.isArray(value)) return value.length ? value.map(item => item && typeof item === 'object' ? `${pad}-\n${yaml(item, indent + 2)}` : `${pad}- ${JSON.stringify(item)}`).join('\n') : `${pad}[]`;
  if (value && typeof value === 'object') return Object.entries(value).map(([key, item]) => {
    if (item && typeof item === 'object' && Object.keys(item).length) return `${pad}${key}:\n${yaml(item, indent + 2)}`;
    return `${pad}${key}: ${JSON.stringify(item)}`;
  }).join('\n');
  return `${pad}${JSON.stringify(value)}`;
}
export function generateFiles(input: OnboardingInput): GeneratedFile[] {
  const errors = validateInput(input);
  if (errors.length) throw new Error(errors.join(' '));
  const id = `${input.name}-${input.environment}`;
  const directory = `minikube/flux/onboarded/${id}`;
  const valuesPath = `minikube/dev/apps/${id}.yaml`;
  const manifest = (path: string, object: unknown): GeneratedFile => ({ path, content: yaml(object) + '\n' });
  const metadata = { name: id, namespace: id };
  return [
    manifest(valuesPath, {
      fullnameOverride: id, replicaCount: input.replicas,
      image: { repository: input.repository, tag: input.tag, pullPolicy: 'IfNotPresent' },
      service: { port: input.port },
      ingress: { enabled: Boolean(input.hostname), className: 'nginx', ...(input.hostname ? { hosts: [{ host: input.hostname, paths: [{ path: '/', pathType: 'Prefix' }] }] } : {}) },
    }),
    manifest(`${directory}/namespace.yaml`, { apiVersion: 'v1', kind: 'Namespace', metadata: { name: id, labels: { 'app.kubernetes.io/name': input.name, 'platform.devops/environment': input.environment } } }),
    manifest(`${directory}/artifact-generator.yaml`, {
      apiVersion: 'source.extensions.fluxcd.io/v1beta1', kind: 'ArtifactGenerator', metadata,
      spec: { sources: [{ alias: 'repo', kind: 'GitRepository', name: 'devops', namespace: 'default' }], artifacts: [{ name: id, copy: [{ from: '@repo/charts/dev/**', to: '@artifact/dev/' }, { from: `@repo/${valuesPath}`, to: '@artifact/dev/values.yaml', strategy: 'Merge' }] }] },
    }),
    manifest(`${directory}/helm-release.yaml`, {
      apiVersion: 'helm.toolkit.fluxcd.io/v2', kind: 'HelmRelease', metadata,
      spec: { interval: '1m', releaseName: id, chartRef: { kind: 'ExternalArtifact', name: id } },
    }),
    manifest(`${directory}/kustomization.yaml`, { apiVersion: 'kustomize.config.k8s.io/v1beta1', kind: 'Kustomization', resources: ['namespace.yaml', 'artifact-generator.yaml', 'helm-release.yaml'] }),
    { path: 'ONBOARDING.md', content: `# Onboard ${id}\n\n1. Extract this ZIP at the Devops repository root.\n2. Add this entry to the resources list in minikube/flux/staging/kustomization.yaml:\n\n\`\`\`yaml\n  - ../onboarded/${id}\n\`\`\`\n\n3. Validate locally: kubectl kustomize minikube/flux/staging\n4. Commit the generated files and resource entry, open a PR and merge into the branch watched by GitRepository/devops (main by default).\n5. Ensure the platform was bootstrapped once with make flux and make flux-app.\n6. Check: flux get helmreleases -n ${id}\n\nFuture upgrades: edit image.tag in ${valuesPath} and merge another PR.\nUse an existing immutable tag from the registry; generating this bundle does not build or publish an image.\n${input.hostname ? `\nFor local kind, route ${input.hostname} to 127.0.0.1 or test with curl -H 'Host: ${input.hostname}' http://localhost/.\n` : `\nFor local access: kubectl -n ${id} port-forward svc/${id} 8081:${input.port}\n`}\nThe shared chart expects a process compatible with its existing command args ([serve]) and port settings. Private images need separately configured imagePullSecrets.\nCross-namespace source references to default/devops must be enabled, as they are in this lab.\nThe archive does not replace existing staging configuration, create a PR or deploy resources. Review generated resource names before overwriting existing files.\n` },
  ];
}

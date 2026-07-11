'use client';

import { FormEvent, useEffect, useState } from 'react';

type SetupStatus = {
  configured: boolean;
  connected: boolean;
  redirect_url?: string;
  missing_config?: string[];
  connected_channels: number;
  scopes: string[];
};

type Channel = {
  id: string;
  title: string;
  thumbnail_url?: string;
  connected_at: string;
};

type PlanResponse = {
  request_id: string;
  model?: string;
  prompt_template?: string;
  summary: string;
  confidence?: number;
  evidence?: string[];
  recommendations?: string[];
  artifacts?: string[];
  timestamp?: string;
};

type UploadResponse = {
  channel_id: string;
  video_id: string;
  upload_status: string;
  privacy_status: string;
  summary: string;
  timestamp: string;
};

type AuthResponse = {
  auth_url: string;
  state: string;
};

type ChannelsResponse = {
  channels: Channel[];
};

const configuredBaseURL = process.env.NEXT_PUBLIC_INFERENCE_BASE_URL ?? 'http://localhost:8080';
const inferenceBaseURL = configuredBaseURL.endsWith('/') ? configuredBaseURL.slice(0, -1) : configuredBaseURL;

async function requestJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(inferenceBaseURL + path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {})
    }
  });

  const text = await response.text();
  const payload = text ? JSON.parse(text) : null;
  if (!response.ok) {
    throw new Error(payload?.error ?? payload?.message ?? 'Request failed');
  }
  return payload as T;
}

function formatTimestamp(value?: string) {
  if (!value) {
    return 'Not available';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
}

export function YouTubeAutomationStudio() {
  const [setupStatus, setSetupStatus] = useState<SetupStatus | null>(null);
  const [channels, setChannels] = useState<Channel[]>([]);
  const [setupError, setSetupError] = useState<string | null>(null);
  const [setupLoading, setSetupLoading] = useState(true);

  const [authResult, setAuthResult] = useState<AuthResponse | null>(null);
  const [authError, setAuthError] = useState<string | null>(null);
  const [authLoading, setAuthLoading] = useState(false);

  const [planPrompt, setPlanPrompt] = useState('Build a YouTube launch brief for a Kubernetes troubleshooting tutorial with title ideas, chapter beats, CTA, and upload notes.');
  const [planResult, setPlanResult] = useState<PlanResponse | null>(null);
  const [planError, setPlanError] = useState<string | null>(null);
  const [planLoading, setPlanLoading] = useState(false);

  const [uploadForm, setUploadForm] = useState({
    channel_id: '',
    file_path: '',
    title: '',
    description: '',
    tags: 'kubernetes, devops, tutorial',
    privacy_status: 'private'
  });
  const [uploadResult, setUploadResult] = useState<UploadResponse | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [uploadLoading, setUploadLoading] = useState(false);

  async function refreshSetup() {
    setSetupLoading(true);
    setSetupError(null);

    try {
      const [status, channelsResponse] = await Promise.all([
        requestJSON<SetupStatus>('/youtube/setup-status'),
        requestJSON<ChannelsResponse>('/youtube/channels')
      ]);

      setSetupStatus(status);
      setChannels(channelsResponse.channels);
      setUploadForm((current) => ({
        ...current,
        channel_id: current.channel_id || channelsResponse.channels[0]?.id || ''
      }));
    } catch (error) {
      setSetupError(error instanceof Error ? error.message : 'Could not load YouTube setup status.');
    } finally {
      setSetupLoading(false);
    }
  }

  useEffect(() => {
    void refreshSetup();
  }, []);

  async function handleStartAuth() {
    setAuthLoading(true);
    setAuthError(null);

    try {
      const result = await requestJSON<AuthResponse>('/youtube/auth/start', { method: 'POST', body: '{}' });
      setAuthResult(result);
    } catch (error) {
      setAuthError(error instanceof Error ? error.message : 'Could not start YouTube OAuth.');
    } finally {
      setAuthLoading(false);
    }
  }

  async function handleGeneratePlan(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setPlanLoading(true);
    setPlanError(null);

    try {
      const result = await requestJSON<PlanResponse>('/generate-video-plan', {
        method: 'POST',
        body: JSON.stringify({ org_id: 'youtube-studio', prompt: planPrompt })
      });
      setPlanResult(result);
    } catch (error) {
      setPlanError(error instanceof Error ? error.message : 'Could not generate the video plan.');
    } finally {
      setPlanLoading(false);
    }
  }

  async function handleUpload(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setUploadLoading(true);
    setUploadError(null);

    try {
      const result = await requestJSON<UploadResponse>('/upload-video', {
        method: 'POST',
        body: JSON.stringify({
          channel_id: uploadForm.channel_id,
          file_path: uploadForm.file_path,
          title: uploadForm.title,
          description: uploadForm.description,
          tags: uploadForm.tags
            .split(',')
            .map((value) => value.trim())
            .filter(Boolean),
          privacy_status: uploadForm.privacy_status
        })
      });
      setUploadResult(result);
    } catch (error) {
      setUploadError(error instanceof Error ? error.message : 'Could not upload the video.');
    } finally {
      setUploadLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      <section className="overflow-hidden rounded-[28px] border border-cyan-400/15 bg-[linear-gradient(135deg,rgba(8,17,31,0.96),rgba(7,30,52,0.94))] shadow-[0_30px_80px_rgba(3,10,24,0.35)]">
        <div className="grid gap-6 px-6 py-8 lg:grid-cols-[1.45fr_0.95fr] lg:px-8">
          <div className="space-y-4">
            <div className="text-xs uppercase tracking-[0.34em] text-cyan-300/70">YouTube Studio Slice</div>
            <div>
              <h2 className="text-3xl font-semibold tracking-tight text-white">Connect a channel, generate a brief, and push a real upload request.</h2>
              <p className="mt-3 max-w-3xl text-sm leading-7 text-slate-300">
                This page talks directly to the inference service on <span className="font-medium text-cyan-200">{inferenceBaseURL}</span>. OAuth setup and uploads stay blocked until Google credentials are configured on the backend.
              </p>
            </div>
            <div className="grid gap-3 sm:grid-cols-3">
              <div className="rounded-2xl border border-white/10 bg-black/20 p-4">
                <div className="text-xs uppercase tracking-[0.28em] text-slate-400">Configured</div>
                <div className="mt-3 text-2xl font-semibold text-white">{setupStatus?.configured ? 'Yes' : 'No'}</div>
              </div>
              <div className="rounded-2xl border border-white/10 bg-black/20 p-4">
                <div className="text-xs uppercase tracking-[0.28em] text-slate-400">Connected Channels</div>
                <div className="mt-3 text-2xl font-semibold text-white">{setupStatus?.connected_channels ?? 0}</div>
              </div>
              <div className="rounded-2xl border border-white/10 bg-black/20 p-4">
                <div className="text-xs uppercase tracking-[0.28em] text-slate-400">Backend Status</div>
                <div className="mt-3 text-2xl font-semibold text-white">{setupLoading ? 'Loading' : setupError ? 'Error' : 'Ready'}</div>
              </div>
            </div>
          </div>
          <div className="rounded-[24px] border border-cyan-300/15 bg-black/25 p-5">
            <div className="flex items-center justify-between gap-3">
              <div>
                <div className="text-xs uppercase tracking-[0.3em] text-cyan-300/70">Live Setup</div>
                <h3 className="mt-2 text-lg font-semibold text-white">Environment and channel state</h3>
              </div>
              <button
                type="button"
                onClick={() => void refreshSetup()}
                className="rounded-full border border-cyan-300/25 px-4 py-2 text-sm text-cyan-100 transition hover:border-cyan-200/50 hover:bg-cyan-300/10"
              >
                Refresh
              </button>
            </div>
            <div className="mt-4 space-y-3 text-sm text-slate-300">
              <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                <div className="font-medium text-white">Redirect URL</div>
                <div className="mt-2 break-all text-slate-300">{setupStatus?.redirect_url ?? 'Not configured'}</div>
              </div>
              <div className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                <div className="font-medium text-white">Required scopes</div>
                <div className="mt-2 flex flex-wrap gap-2">
                  {(setupStatus?.scopes ?? []).map((scope) => (
                    <span key={scope} className="rounded-full border border-cyan-400/20 bg-cyan-400/10 px-3 py-1 text-xs text-cyan-100">
                      {scope}
                    </span>
                  ))}
                </div>
              </div>
              {!setupLoading && setupStatus?.missing_config?.length ? (
                <div className="rounded-2xl border border-amber-400/25 bg-amber-400/10 p-4 text-amber-100">
                  Missing backend env: {setupStatus.missing_config.join(', ')}
                </div>
              ) : null}
              {setupError ? <div className="rounded-2xl border border-rose-400/25 bg-rose-400/10 p-4 text-rose-100">{setupError}</div> : null}
            </div>
          </div>
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[1.05fr_1fr]">
        <section className="space-y-6">
          <div className="rounded-[24px] border border-white/10 bg-black/20 p-5">
            <div className="flex items-center justify-between gap-4">
              <div>
                <div className="text-xs uppercase tracking-[0.3em] text-slate-400">OAuth</div>
                <h3 className="mt-2 text-xl font-semibold text-white">Connect a YouTube channel</h3>
              </div>
              <button
                type="button"
                onClick={() => void handleStartAuth()}
                disabled={authLoading}
                className="rounded-full bg-cyan-300 px-4 py-2 text-sm font-medium text-slate-950 transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-100/40"
              >
                {authLoading ? 'Creating URL...' : 'Start OAuth'}
              </button>
            </div>
            <p className="mt-3 text-sm leading-6 text-slate-300">
              The backend generates a Google OAuth URL with upload and read-only scopes, then stores the state token in memory until the callback lands on the inference service.
            </p>
            {authError ? <div className="mt-4 rounded-2xl border border-rose-400/25 bg-rose-400/10 p-4 text-sm text-rose-100">{authError}</div> : null}
            {authResult ? (
              <div className="mt-4 rounded-2xl border border-cyan-300/20 bg-cyan-400/10 p-4 text-sm text-cyan-50">
                <div className="font-medium">Auth URL ready</div>
                <a href={authResult.auth_url} target="_blank" rel="noreferrer" className="mt-2 block break-all text-cyan-100 underline decoration-cyan-200/60 underline-offset-4">
                  {authResult.auth_url}
                </a>
                <div className="mt-3 text-xs uppercase tracking-[0.28em] text-cyan-200/70">State {authResult.state}</div>
              </div>
            ) : null}
          </div>

          <div className="rounded-[24px] border border-white/10 bg-black/20 p-5">
            <div className="text-xs uppercase tracking-[0.3em] text-slate-400">Connected Channels</div>
            <h3 className="mt-2 text-xl font-semibold text-white">In-memory channel cache</h3>
            <div className="mt-4 space-y-3">
              {channels.length === 0 ? (
                <div className="rounded-2xl border border-dashed border-white/10 bg-slate-950/30 p-4 text-sm text-slate-300">
                  No channels are connected yet. Complete the OAuth callback with valid Google credentials to populate this list.
                </div>
              ) : (
                channels.map((channel) => (
                  <div key={channel.id} className="flex items-center gap-4 rounded-2xl border border-white/10 bg-slate-950/35 p-4">
                    <div className="flex h-12 w-12 items-center justify-center overflow-hidden rounded-2xl border border-white/10 bg-slate-900 text-xs text-slate-400">
                      {channel.thumbnail_url ? <img src={channel.thumbnail_url} alt={channel.title} className="h-full w-full object-cover" /> : 'YT'}
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="truncate font-medium text-white">{channel.title}</div>
                      <div className="truncate text-sm text-slate-400">{channel.id}</div>
                    </div>
                    <div className="text-xs uppercase tracking-[0.26em] text-slate-500">{formatTimestamp(channel.connected_at)}</div>
                  </div>
                ))
              )}
            </div>
          </div>
        </section>

        <section className="space-y-6">
          <form onSubmit={handleGeneratePlan} className="rounded-[24px] border border-white/10 bg-black/20 p-5">
            <div className="text-xs uppercase tracking-[0.3em] text-slate-400">Creative Brief</div>
            <h3 className="mt-2 text-xl font-semibold text-white">Generate a video plan</h3>
            <textarea
              value={planPrompt}
              onChange={(event) => setPlanPrompt(event.target.value)}
              rows={6}
              className="mt-4 w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
              placeholder="Describe the topic, audience, angle, and publishing goal."
            />
            <div className="mt-4 flex items-center justify-between gap-3">
              <p className="text-sm text-slate-400">Posts to <span className="text-cyan-200">/generate-video-plan</span> and returns the routed model template.</p>
              <button
                type="submit"
                disabled={planLoading}
                className="rounded-full bg-white px-4 py-2 text-sm font-medium text-slate-950 transition hover:bg-slate-200 disabled:cursor-not-allowed disabled:bg-slate-200/50"
              >
                {planLoading ? 'Generating...' : 'Generate Plan'}
              </button>
            </div>
            {planError ? <div className="mt-4 rounded-2xl border border-rose-400/25 bg-rose-400/10 p-4 text-sm text-rose-100">{planError}</div> : null}
            {planResult ? (
              <div className="mt-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4 text-sm text-slate-200">
                <div className="flex flex-wrap gap-3 text-xs uppercase tracking-[0.25em] text-slate-400">
                  <span>Model {planResult.model ?? 'n/a'}</span>
                  <span>Template {planResult.prompt_template ?? 'n/a'}</span>
                  <span>Confidence {planResult.confidence ?? 0}</span>
                </div>
                <p className="mt-4 leading-7 text-slate-100">{planResult.summary}</p>
                <pre className="mt-4 overflow-x-auto rounded-2xl border border-white/10 bg-black/30 p-4 text-xs text-cyan-100">{JSON.stringify(planResult, null, 2)}</pre>
              </div>
            ) : null}
          </form>

          <form onSubmit={handleUpload} className="rounded-[24px] border border-white/10 bg-black/20 p-5">
            <div className="text-xs uppercase tracking-[0.3em] text-slate-400">Upload Request</div>
            <h3 className="mt-2 text-xl font-semibold text-white">Submit a real YouTube upload</h3>
            <div className="mt-4 grid gap-4 md:grid-cols-2">
              <label className="block text-sm text-slate-300">
                <span className="mb-2 block">Channel</span>
                <select
                  value={uploadForm.channel_id}
                  onChange={(event) => setUploadForm((current) => ({ ...current, channel_id: event.target.value }))}
                  className="w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
                >
                  <option value="">Select a connected channel</option>
                  {channels.map((channel) => (
                    <option key={channel.id} value={channel.id}>
                      {channel.title} ({channel.id})
                    </option>
                  ))}
                </select>
              </label>
              <label className="block text-sm text-slate-300">
                <span className="mb-2 block">Privacy</span>
                <select
                  value={uploadForm.privacy_status}
                  onChange={(event) => setUploadForm((current) => ({ ...current, privacy_status: event.target.value }))}
                  className="w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
                >
                  <option value="private">private</option>
                  <option value="unlisted">unlisted</option>
                  <option value="public">public</option>
                </select>
              </label>
              <label className="block text-sm text-slate-300 md:col-span-2">
                <span className="mb-2 block">Local file path</span>
                <input
                  value={uploadForm.file_path}
                  onChange={(event) => setUploadForm((current) => ({ ...current, file_path: event.target.value }))}
                  className="w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
                  placeholder="/absolute/path/to/video.mp4"
                />
              </label>
              <label className="block text-sm text-slate-300 md:col-span-2">
                <span className="mb-2 block">Title</span>
                <input
                  value={uploadForm.title}
                  onChange={(event) => setUploadForm((current) => ({ ...current, title: event.target.value }))}
                  className="w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
                  placeholder="Kubernetes debugging in 7 minutes"
                />
              </label>
              <label className="block text-sm text-slate-300 md:col-span-2">
                <span className="mb-2 block">Description</span>
                <textarea
                  value={uploadForm.description}
                  onChange={(event) => setUploadForm((current) => ({ ...current, description: event.target.value }))}
                  rows={4}
                  className="w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
                  placeholder="Add the publish-ready description used for the YouTube upload."
                />
              </label>
              <label className="block text-sm text-slate-300 md:col-span-2">
                <span className="mb-2 block">Tags</span>
                <input
                  value={uploadForm.tags}
                  onChange={(event) => setUploadForm((current) => ({ ...current, tags: event.target.value }))}
                  className="w-full rounded-2xl border border-white/10 bg-slate-950/40 px-4 py-3 text-sm text-slate-100 outline-none transition focus:border-cyan-300/40"
                  placeholder="kubernetes, devops, tutorial"
                />
              </label>
            </div>
            <div className="mt-4 flex items-center justify-between gap-3">
              <p className="text-sm text-slate-400">Posts to <span className="text-cyan-200">/upload-video</span> with channel, file path, metadata, tags, and privacy.</p>
              <button
                type="submit"
                disabled={uploadLoading}
                className="rounded-full bg-cyan-300 px-4 py-2 text-sm font-medium text-slate-950 transition hover:bg-cyan-200 disabled:cursor-not-allowed disabled:bg-cyan-100/40"
              >
                {uploadLoading ? 'Uploading...' : 'Upload Video'}
              </button>
            </div>
            {uploadError ? <div className="mt-4 rounded-2xl border border-rose-400/25 bg-rose-400/10 p-4 text-sm text-rose-100">{uploadError}</div> : null}
            {uploadResult ? (
              <div className="mt-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4 text-sm text-slate-200">
                <div className="grid gap-3 sm:grid-cols-2">
                  <div>
                    <div className="text-xs uppercase tracking-[0.25em] text-slate-400">Video ID</div>
                    <div className="mt-2 break-all font-medium text-white">{uploadResult.video_id}</div>
                  </div>
                  <div>
                    <div className="text-xs uppercase tracking-[0.25em] text-slate-400">Upload Status</div>
                    <div className="mt-2 font-medium text-white">{uploadResult.upload_status}</div>
                  </div>
                </div>
                <div className="mt-4 text-slate-100">{uploadResult.summary}</div>
                <div className="mt-2 text-xs uppercase tracking-[0.24em] text-slate-500">{formatTimestamp(uploadResult.timestamp)}</div>
                <pre className="mt-4 overflow-x-auto rounded-2xl border border-white/10 bg-black/30 p-4 text-xs text-cyan-100">{JSON.stringify(uploadResult, null, 2)}</pre>
              </div>
            ) : null}
          </form>
        </section>
      </div>
    </div>
  );
}

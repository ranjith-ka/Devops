export type InferenceContext = {
  branch?: string;
  environment?: string;
  timeRange?: string;
};

export type InferenceRequest = {
  orgId: string;
  repository?: string;
  pipeline?: string;
  deployment?: string;
  prompt: string;
  context?: InferenceContext;
};

export type InferenceResponse = {
  requestId: string;
  model?: string;
  promptTemplate?: string;
  summary: string;
  confidence?: number;
  evidence?: string[];
  recommendations?: string[];
  artifacts?: string[];
};

export type EventEnvelope<TPayload = Record<string, unknown>> = {
  eventId: string;
  eventType: string;
  schemaVersion: number;
  tenantId: string;
  correlationId: string;
  causationId?: string;
  occurredAt: string;
  source: string;
  payload: TPayload;
};

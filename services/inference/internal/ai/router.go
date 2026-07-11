package ai

import "strings"

type ProviderPlan struct {
  Provider       string  `json:"provider"`
  Model          string  `json:"model"`
  Fallback       string  `json:"fallback"`
  PromptTemplate string  `json:"prompt_template"`
  Reason         string  `json:"reason"`
  Confidence     float64 `json:"confidence"`
}

type Router struct {
  catalog *Catalog
}

func NewRouter() *Router { return &Router{catalog: NewCatalog()} }

func (r *Router) Plan(task string, prompt string) ProviderPlan {
  profile := r.catalog.Resolve(task, prompt)
  return ProviderPlan{
    Provider:       profile.Provider,
    Model:          profile.Model,
    Fallback:       profile.Fallback,
    PromptTemplate: profile.Template,
    Reason:         profile.Reason,
    Confidence:     profile.Confidence,
  }
}

type Catalog struct {
  profiles       []TaskProfile
  defaultProfile TaskProfile
}

type TaskProfile struct {
  Matchers   []string
  Provider   string
  Model      string
  Fallback   string
  Template   string
  Reason     string
  Confidence float64
}

func NewCatalog() *Catalog {
  return &Catalog{
    profiles: []TaskProfile{
      {
        Matchers:   []string{"security", "policy", "vulnerability", "scan"},
        Provider:   "claude",
        Model:      "claude-3.5-sonnet",
        Fallback:   "gpt",
        Template:   "security-review-v1",
        Reason:     "Long-context policy analysis and explanation quality.",
        Confidence: 0.91,
      },
      {
        Matchers:   []string{"terraform", "yaml", "kubernetes", "helm"},
        Provider:   "gpt",
        Model:      "gpt-5.4-mini",
        Fallback:   "qwen",
        Template:   "delivery-artifact-v1",
        Reason:     "High-quality structured generation for infra artifacts.",
        Confidence: 0.93,
      },
      {
        Matchers:   []string{"cost", "scale", "runner", "capacity", "utilization"},
        Provider:   "qwen",
        Model:      "qwen2.5-coder",
        Fallback:   "llama",
        Template:   "optimization-plan-v1",
        Reason:     "Efficient reasoning for optimization and self-hosted environments.",
        Confidence: 0.88,
      },
      {
        Matchers:   []string{"log", "incident", "summary", "timeline", "deploy"},
        Provider:   "claude",
        Model:      "claude-3.5-sonnet",
        Fallback:   "gpt",
        Template:   "incident-summary-v1",
        Reason:     "Long-context summarization and incident reconstruction.",
        Confidence: 0.9,
      },
      {
        Matchers:   []string{"youtube", "video", "thumbnail", "shorts", "upload"},
        Provider:   "gpt",
        Model:      "gpt-5.4-mini",
        Fallback:   "claude",
        Template:   "youtube-automation-v1",
        Reason:     "Structured creative planning and metadata generation for repeatable publishing workflows.",
        Confidence: 0.92,
      },
    },
    defaultProfile: TaskProfile{
      Provider:   "gpt",
      Model:      "gpt-5.4-mini",
      Fallback:   "deepseek",
      Template:   "general-assistant-v1",
      Reason:     "Balanced reasoning, tool calling, and routing quality.",
      Confidence: 0.89,
    },
  }
}

func (c *Catalog) Resolve(task string, prompt string) TaskProfile {
  normalized := strings.ToLower(task + " " + prompt)
  for _, profile := range c.profiles {
    for _, matcher := range profile.Matchers {
      if strings.Contains(normalized, matcher) {
        return profile
      }
    }
  }
  return c.defaultProfile
}

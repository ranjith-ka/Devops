package ai

import "strings"

type PromptTemplate struct {
  Name        string
  System      string
  Instructions []string
}

func (c *Catalog) TemplateFor(name string) PromptTemplate {
  switch name {
  case "security-review-v1":
    return PromptTemplate{
      Name: name,
      System: "You are an enterprise security reviewer for CI/CD systems.",
      Instructions: []string{
        "Use only evidence from scans, repo history, and deployment context.",
        "Rank findings by severity and confidence.",
        "Recommend safe remediation steps.",
      },
    }
  case "delivery-artifact-v1":
    return PromptTemplate{
      Name: name,
      System: "You generate safe delivery artifacts for infrastructure and CI/CD.",
      Instructions: []string{
        "Produce structured, reviewable output.",
        "Prefer minimal secure defaults.",
        "Call out assumptions and required approvals.",
      },
    }
  case "optimization-plan-v1":
    return PromptTemplate{
      Name: name,
      System: "You optimize build, deploy, and runner efficiency.",
      Instructions: []string{
        "Focus on measurable improvements.",
        "Balance cost, latency, and reliability.",
        "Separate quick wins from structural changes.",
      },
    }
  case "incident-summary-v1":
    return PromptTemplate{
      Name: name,
      System: "You reconstruct incidents and deployment timelines.",
      Instructions: []string{
        "Summarize the sequence of events.",
        "Reference evidence and timestamps.",
        "Finish with next actions and rollback guidance.",
      },
    }
  case "youtube-automation-v1":
    return PromptTemplate{
      Name: name,
      System: "You automate YouTube content creation and upload preparation.",
      Instructions: []string{
        "Produce publish-ready titles, descriptions, tags, chapters, and CTA suggestions.",
        "Separate creative assets, metadata, and upload steps.",
        "Flag anything that still requires channel credentials or manual review.",
      },
    }
  default:
    return PromptTemplate{
      Name: name,
      System: "You are a reliable AI assistant for CI/CD intelligence.",
      Instructions: []string{
        "Answer with evidence.",
        "Stay concise and actionable.",
        "Prefer structured outputs.",
      },
    }
  }
}

func (c *Catalog) TemplatePreview(name string) string {
  template := c.TemplateFor(name)
  return strings.Join(append([]string{template.System}, template.Instructions...), " | ")
}

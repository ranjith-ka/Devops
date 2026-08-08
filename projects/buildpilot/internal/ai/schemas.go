package ai

import "encoding/json"

var BuildPlanSchema = json.RawMessage(`{
  "type": "object",
  "properties": {
    "summary": {"type": "string"},
    "language": {"type": "string"},
    "framework": {"type": "string"},
    "builder": {"type": "string"},
    "platforms": {"type": "array", "items": {"type": "string"}},
    "container_port": {"type": "integer"},
    "health_endpoint": {"type": "string"},
    "required_files": {"type": "array", "items": {"type": "string"}},
    "warnings": {"type": "array", "items": {"type": "string"}},
    "needs_user_approval": {"type": "boolean"}
  },
  "required": ["summary", "language", "framework", "builder", "platforms", "container_port", "health_endpoint", "required_files", "warnings", "needs_user_approval"]
}`)

var BuildDiagnosisSchema = json.RawMessage(`{
  "type": "object",
  "properties": {
    "summary": {"type": "string"},
    "category": {"type": "string"},
    "confidence": {"type": "number", "minimum": 0, "maximum": 1},
    "evidence": {"type": "array", "items": {"type": "string"}},
    "suggested_changes": {"type": "array", "items": {"type": "string"}},
    "safe_to_auto_apply": {"type": "boolean"}
  },
  "required": ["summary", "category", "confidence", "evidence", "suggested_changes", "safe_to_auto_apply"]
}`)

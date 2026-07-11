import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function AiAssistantPage() {
  return (
    <AppShell activePath="/ai-assistant">
      <ScreenSpec
        title="AI Assistant"
        subtitle="A context-rich assistant that reasons over delivery evidence and returns ranked, actionable guidance rather than generic chat replies."
        api={["POST /chat", "POST /root-cause", "POST /optimize", "POST /security-review"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Context panel                        | Conversational reasoning                  |
| repo, pipeline, release, logs        | evidence-backed answers and actions       |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Context rail > evidence bundle > thread", "Assistant output > citations > action chips", "Approval flow for generated automation"]}
        journey={["Ask a system question in natural language.", "Provide evidence from the selected repository or release.", "Receive a plan, summary, and follow-up actions."]}
        stateModel={["Thread state, selected evidence, memory summary, and policy mode shape the response.", "Assistant memory is tenant-scoped and versioned."]}
        states={[
          { label: 'Loading', description: 'Show a typing indicator, cached evidence, and an active context badge.' },
          { label: 'Error', description: 'Render provider fallback status with a retry path and cached previous answer.' },
          { label: 'Empty', description: 'Prompt the user to select a repository, pipeline, or release context.' }
        ]}
        evidence={["logs", "metrics", "PR diffs", "deployment history", "security findings"]}
      />
    </AppShell>
  );
}

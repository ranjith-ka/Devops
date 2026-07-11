import { AppShell } from '@/components/app-shell';
import { ScreenSpec } from '@/components/screen-spec';

export default function LogsViewerPage() {
  return (
    <AppShell activePath="/logs-viewer">
      <ScreenSpec
        title="Logs Viewer"
        subtitle="Search, filter, and explain logs with AI-assisted evidence selection and highlighted failure correlations."
        api={["POST /chat", "POST /root-cause", "POST /deployment-summary"]}
        wireframe={`+--------------------------------------+------------------------------------------+
| Search and filters                   | Highlighted log stream                   |
| source, time, severity               | errors, correlations, inline evidence     |
+--------------------------------------+------------------------------------------+`}
        hierarchy={["Search bar > filters > live stream", "Highlighted lines > related traces > evidence drawer", "AI explanation panel and action chips"]}
        journey={["Query logs by service, release, or failure signature.", "Open highlighted spans and follow correlated evidence.", "Ask the assistant to explain the likely cause."]}
        stateModel={["Time range, service, trace correlation, and log query state drive the stream.", "Highlighted evidence updates as the user moves the log cursor."]}
        states={[
          { label: 'Loading', description: 'Show log skeleton rows and streaming placeholder text.' },
          { label: 'Error', description: 'Indicate source unavailability, query timeout, or parsing failure.' },
          { label: 'Empty', description: 'Explain that no logs match the query and suggest broader filters.' }
        ]}
        evidence={["stdout", "error traces", "trace ids", "deployment metadata"]}
      />
    </AppShell>
  );
}

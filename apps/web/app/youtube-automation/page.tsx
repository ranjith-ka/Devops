import { AppShell } from '@/components/app-shell';
import { YouTubeAutomationStudio } from '@/components/youtube-automation-studio';

export default function YouTubeAutomationPage() {
  return (
    <AppShell activePath="/youtube-automation">
      <YouTubeAutomationStudio />
    </AppShell>
  );
}

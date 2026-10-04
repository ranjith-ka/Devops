import type { Metadata } from 'next';
import { PlatformOnboarding } from '@/components/platform-onboarding';
export const metadata: Metadata = { title: 'Application onboarding | Pipeline OS', description: 'Prepare an application deployment for your Flux platform.' };
export default function Page() { return <PlatformOnboarding />; }

import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geist = Geist({ subsets: ["latin"], variable: "--font-sans" });
const mono = Geist_Mono({ subsets: ["latin"], variable: "--font-mono" });

export const metadata: Metadata = {
  title: "Ranjith K A | Platform, DevOps & AI Infrastructure",
  description:
    "Senior platform and DevOps engineer helping teams build reliable Kubernetes platforms, delivery systems, observability and production AI infrastructure.",
  openGraph: {
    title: "Ranjith K A | Platform, DevOps & AI Infrastructure",
    description:
      "Platform engineering, Kubernetes, developer enablement, observability and AI infrastructure consulting.",
    type: "website",
  },
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className={`${geist.variable} ${mono.variable}`}>{children}</body>
    </html>
  );
}

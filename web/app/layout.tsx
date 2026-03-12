import type { Metadata } from "next";
import { SessionProvider } from "next-auth/react";
import FeedbackWidget from "@/components/FeedbackWidget";
import "./globals.css";

export const metadata: Metadata = {
  title: "Lark - Learn Languages Through Adventure",
  description:
    "A text-adventure game that teaches you real-world language skills through immersive role-play scenarios. 80+ languages, 40 scenarios, grammar correction.",
  openGraph: {
    title: "Lark - Learn Languages Through Adventure",
    description:
      "A text-adventure game that teaches you real-world language skills through immersive role-play scenarios. 80+ languages, 40 scenarios, grammar correction.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <head>
        <script
          async
          src="https://www.googletagmanager.com/gtag/js?id=G-VTSTFQZH2D"
        />
        <script
          dangerouslySetInnerHTML={{
            __html: `
              window.dataLayer = window.dataLayer || [];
              function gtag(){dataLayer.push(arguments);}
              gtag('js', new Date());
              gtag('config', 'G-VTSTFQZH2D');
            `,
          }}
        />
        <link
          href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="antialiased">
        <SessionProvider>
          {children}
          <FeedbackWidget />
        </SessionProvider>
      </body>
    </html>
  );
}

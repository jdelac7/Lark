"use client";

import CopyCmd from "@/components/CopyCmd";

export default function SuccessContent({
  licenseKey,
  error,
}: {
  licenseKey: string | null;
  error: string | null;
}) {
  return (
    <div className="flex min-h-screen items-center justify-center px-6">
      <div className="w-full max-w-lg">
        <div className="rounded-lg border border-accent/30 bg-bg-card p-8">
          <div className="text-center">
            <div className="text-4xl text-green">&#10003;</div>
            <h1 className="mt-4 text-2xl font-bold">Welcome to Lark!</h1>
            <p className="mt-2 text-text-dim">
              Your subscription is active. Here&apos;s how to get started:
            </p>
          </div>

          {error && (
            <div className="mt-6 rounded-md border border-yellow/30 bg-yellow/10 p-4 text-sm text-yellow">
              {error}
            </div>
          )}

          <div className="mt-8 space-y-4">
            <h2 className="font-semibold">Activation Instructions</h2>
            <div className="space-y-3 text-sm">
              <div className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-bold text-accent">
                  1
                </span>
                <div className="min-w-0 flex-1">
                  <p>Install Lark (if you haven&apos;t already):</p>
                  <div className="mt-1 rounded border border-border bg-bg px-3 py-2">
                    <CopyCmd cmd="curl -fsSL https://lark.black/install.sh | sh" className="text-green" />
                  </div>
                </div>
              </div>
              <div className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-bold text-accent">
                  2
                </span>
                <div>
                  <p>Activate your license:</p>
                  <code className="mt-1 block rounded border border-border bg-bg px-3 py-2 text-green">
                    lark activate{" "}
                    {licenseKey ? licenseKey : "<YOUR-LICENSE-KEY>"}
                  </code>
                </div>
              </div>
              <div className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-bold text-accent">
                  3
                </span>
                <div>
                  <p>Start playing:</p>
                  <code className="mt-1 block rounded border border-border bg-bg px-3 py-2 text-green">
                    lark
                  </code>
                </div>
              </div>
            </div>
          </div>

          {/* Browser play prompt */}
          <div className="mt-8 border-t border-border pt-6">
            <div className="mb-2 text-sm font-semibold text-cyan">
              Play in the browser
            </div>
            <p className="mb-4 text-sm text-text-dim">
              Your subscription is linked to your account. You can play directly
              in the browser — no install needed.
            </p>
            <a
              href="/play"
              className="block w-full border border-cyan py-3 text-center text-sm font-bold text-cyan transition-colors hover:bg-cyan/10"
            >
              &gt; PLAY NOW
            </a>
          </div>

          <div className="mt-6 text-center">
            <a
              href="/"
              className="text-sm text-accent transition-colors hover:text-accent-dim"
            >
              &larr; Back to home
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

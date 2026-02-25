"use client";

import { useSession } from "next-auth/react";
import { useEffect, useState } from "react";
import CopyCmd from "./CopyCmd";

export default function HowItWorks() {
  const { data: session } = useSession();
  const isSubscribed = session?.user?.subscribed;
  const [licenseKey, setLicenseKey] = useState<string | null>(null);

  useEffect(() => {
    if (!isSubscribed) return;
    fetch("/api/license-key")
      .then((r) => r.json())
      .then((d) => setLicenseKey(d.key))
      .catch(() => {});
  }, [isSubscribed]);

  const keyDisplay = licenseKey || "YOUR-LICENSE-KEY";

  return (
    <section id="quickstart" className="border-t border-border py-20">
      <div className="mx-auto max-w-4xl px-6">
        {/* Command prompt */}
        <div className="mb-8 text-sm">
          <span className="text-accent">user@lark</span>
          <span className="text-text-dim">:</span>
          <span className="text-cyan">~</span>
          <span className="text-text-dim">$ </span>
          <span className="text-text">lark quickstart</span>
        </div>

        <div className="mb-8 text-sm font-bold text-yellow">
          QUICKSTART GUIDE
        </div>

        {/* Step 1 — shared for CLI paths */}
        <div className="mb-10 text-sm">
          <div className="mb-2 text-text-dim">
            <span className="text-accent">Step 1</span> &mdash; Install{" "}
            <span className="text-text-dim">(CLI only &mdash; skip for browser)</span>
          </div>
          <div className="border border-border bg-bg-secondary p-4">
            <CopyCmd cmd="curl -fsSL https://lark.black/install.sh | sh" className="text-green" />
            <div className="mt-2 text-accent">
              ✓ Lark installed to ~/.local/bin/lark
            </div>
          </div>
        </div>

        <div className="mb-6 text-sm font-bold text-yellow">
          CHOOSE YOUR PATH
        </div>

        <div className="grid grid-cols-1 gap-x-8 gap-y-5 text-sm md:grid-cols-3 md:grid-rows-[auto_auto_1fr]">
          {/* Path A — BYOK */}
          <div className="grid grid-rows-subgrid gap-y-5 md:row-span-3">
            <div className="text-sm font-bold text-cyan">
              FREE &mdash; BRING YOUR OWN KEY
            </div>

            <div>
              <div className="mb-2 text-text-dim">
                <span className="text-cyan">Step 2</span> &mdash; Get an API
                key
              </div>
              <div className="border border-border bg-bg-secondary p-4">
                <div className="text-green">
                  $ open https://openrouter.ai/keys
                </div>
                <div className="mt-2 text-text-dim">
                  &rarr; Create a free account on OpenRouter
                </div>
                <div className="text-text-dim">
                  &rarr; Generate an API key
                </div>
              </div>
            </div>

            <div>
              <div className="mb-2 text-text-dim">
                <span className="text-cyan">Step 3</span> &mdash; Set your key
                &amp; play
              </div>
              <div className="border border-border bg-bg-secondary p-4">
                <div className="text-green">
                  $ lark apikey YOUR-OPENROUTER-KEY
                </div>
                <div className="mt-2 text-cyan">✓ API key saved!</div>
                <div className="mt-2 text-green">$ lark</div>
                <div className="mt-2 text-cyan">
                  &nbsp;&nbsp;Welcome back, adventurer. Where shall we go today?
                </div>
              </div>
            </div>
          </div>

          {/* Path B — Subscription (CLI) */}
          <div className="grid grid-rows-subgrid gap-y-5 md:row-span-3">
            <div className="text-sm font-bold text-accent">
              SUBSCRIPTION &mdash; CLI
            </div>

            <div>
              <div className="mb-2 text-text-dim">
                <span className="text-accent">Step 2</span> &mdash; Subscribe
              </div>
              <div className="border border-border bg-bg-secondary p-4">
                <div className="text-green">
                  $ open https://lark.black
                </div>
                <div className="mt-2 text-text-dim">
                  &rarr; Subscribe for $2.99/month
                </div>
                <div className="text-text-dim">
                  &rarr; Receive your license key instantly
                </div>
              </div>
            </div>

            <div>
              <div className="mb-2 text-text-dim">
                <span className="text-accent">Step 3</span> &mdash; Activate
                &amp; play
              </div>
              <div className="border border-border bg-bg-secondary p-4">
                <div className="text-green">
                  $ lark activate {keyDisplay}
                </div>
                <div className="mt-2 text-accent">
                  ✓ License activated successfully
                </div>
                <div className="mt-2 text-green">$ lark</div>
                <div className="mt-2 text-accent">
                  &nbsp;&nbsp;Welcome back, adventurer. Where shall we go today?
                </div>
              </div>
            </div>
          </div>

          {/* Path C — Play in Browser */}
          <div className="grid grid-rows-subgrid gap-y-5 md:row-span-3">
            <div className="text-sm font-bold text-purple">
              SUBSCRIPTION &mdash; BROWSER
            </div>

            <div>
              <div className="mb-2 text-text-dim">
                <span className="text-purple">Step 2</span> &mdash; Register
                &amp; subscribe
              </div>
              <div className="border border-border bg-bg-secondary p-4">
                <div className="text-text-dim">
                  &rarr;{" "}
                  <a
                    href="/register"
                    className="text-purple underline decoration-purple/40 hover:decoration-purple"
                  >
                    Create an account
                  </a>
                </div>
                <div className="text-text-dim">
                  &rarr; Subscribe for $2.99/month
                </div>
                <div className="mt-2 text-purple">
                  ✓ Subscription linked automatically
                </div>
              </div>
            </div>

            <div>
              <div className="mb-2 text-text-dim">
                <span className="text-purple">Step 3</span> &mdash; Play
              </div>
              <div className="border border-border bg-bg-secondary p-4">
                <div className="text-text-dim">
                  &rarr; No install needed
                </div>
                <div className="mt-2 text-purple">
                  ✓ Play directly in your browser
                </div>
                <div className="mt-3">
                  <a
                    href="/play"
                    className="inline-block border border-purple px-4 py-2 text-sm font-bold text-purple transition-colors hover:bg-purple/10"
                  >
                    &gt; PLAY NOW
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

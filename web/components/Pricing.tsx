"use client";

import { useSession } from "next-auth/react";
import { useEffect, useState } from "react";
import CopyCmd from "./CopyCmd";

const sharedFeatures = [
  "All 40 built-in scenarios",
  "80+ languages",
  "Unlimited custom scenarios",
  "Instant grammar correction",
  "Vocabulary tracking",
  "All future updates",
];

export default function Pricing() {
  const { data: session } = useSession();
  const isSubscribed = session?.user?.subscribed;
  const productId = process.env.NEXT_PUBLIC_POLAR_PRODUCT_ID;
  const checkoutUrl = `/api/checkout?products=${productId}`;
  const [licenseKey, setLicenseKey] = useState<string | null>(null);

  useEffect(() => {
    if (!isSubscribed) return;
    fetch("/api/license-key")
      .then((r) => r.json())
      .then((d) => setLicenseKey(d.key))
      .catch(() => {});
  }, [isSubscribed]);

  // If not logged in, send to register first — checkout requires an account
  const subscribeHref = session ? checkoutUrl : `/register?callbackUrl=${encodeURIComponent(checkoutUrl)}`;

  return (
    <section id="shop" className="border-t border-border py-20">
      <div className="mx-auto max-w-4xl px-6">
        {/* Command prompt */}
        <div className="mb-8 text-sm">
          <span className="text-accent">user@lark</span>
          <span className="text-text-dim">:</span>
          <span className="text-cyan">~</span>
          <span className="text-text-dim">$ </span>
          <span className="text-text">lark shop</span>
        </div>

        <div className="mb-4 text-sm text-yellow">
          A merchant approaches you...
        </div>

        <div className="mb-8 text-sm italic text-yellow">
          &ldquo;Greetings, adventurer. I have two offers for you today.&rdquo;
        </div>

        <div className="mx-auto grid max-w-3xl grid-cols-1 gap-6 md:grid-cols-2">
          {/* Free tier — BYOK */}
          <div className="border border-border bg-bg-secondary p-6 sm:p-8">
            <div className="mb-1 text-lg font-bold text-cyan">
              BRING YOUR OWN KEY
            </div>
            <div className="mb-6 flex items-baseline gap-1">
              <span className="text-3xl font-bold text-text">Free</span>
            </div>

            <div className="mb-8 space-y-2 text-sm">
              {sharedFeatures.map((feature) => (
                <div key={feature} className="flex items-center gap-2">
                  <span className="text-cyan">✓</span>
                  <span className="text-text">{feature}</span>
                </div>
              ))}
              <div className="flex items-center gap-2">
                <span className="text-text-dim">*</span>
                <span className="text-text-dim">
                  Requires{" "}
                  <a
                    href="https://openrouter.ai/keys"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-cyan underline decoration-cyan/40 hover:decoration-cyan"
                  >
                    OpenRouter
                  </a>{" "}
                  API key
                </span>
              </div>
            </div>

            <a
              href="#quickstart"
              className="block w-full border border-cyan py-3 text-center text-sm font-bold text-cyan transition-colors hover:bg-cyan/10"
            >
              &gt; GET STARTED
            </a>
          </div>

          {/* Pro tier — Subscription */}
          <div className="border border-accent bg-bg-secondary p-6 sm:p-8">
            <div className="mb-1 text-lg font-bold text-accent">
              LARK LICENSE
            </div>
            <div className="mb-6 flex items-baseline gap-1">
              <span className="text-3xl font-bold text-text">$2.99</span>
              <span className="text-sm text-text-dim">/month</span>
            </div>

            <div className="mb-8 space-y-2 text-sm">
              {sharedFeatures.map((feature) => (
                <div key={feature} className="flex items-center gap-2">
                  <span className="text-accent">✓</span>
                  <span className="text-text">{feature}</span>
                </div>
              ))}
              <div className="flex items-center gap-2">
                <span className="text-accent">✓</span>
                <span className="text-text">No API key needed</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-accent">✓</span>
                <span className="text-text">Play in browser</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-accent">✓</span>
                <span className="text-text">Supports development</span>
              </div>
            </div>

            {isSubscribed ? (
              <>
                <div className="mb-3 border border-accent/30 bg-accent/10 py-3 text-center text-sm font-bold text-accent">
                  ✓ SUBSCRIBED
                </div>
                {licenseKey && (
                  <div className="mb-3 border border-border bg-bg p-3 text-xs">
                    <CopyCmd cmd={`lark activate ${licenseKey}`} className="text-accent" />
                  </div>
                )}
                <a
                  href="https://polar.sh/settings/subscriptions"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block w-full border border-text-dim py-2 text-center text-xs text-text-dim transition-colors hover:border-yellow hover:text-yellow"
                >
                  &gt; MANAGE SUBSCRIPTION
                </a>
              </>
            ) : (
              <>
                <a
                  href={subscribeHref}
                  className="block w-full border border-accent py-3 text-center text-sm font-bold text-accent transition-colors hover:bg-accent/10"
                >
                  &gt; {session ? "ACCEPT OFFER" : "CREATE ACCOUNT & SUBSCRIBE"}
                </a>
                <p className="mt-4 text-center text-xs text-text-dim">
                  {!session && "Account required · "}Cancel anytime &middot; Powered by Polar
                </p>
              </>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}

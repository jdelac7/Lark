const features = [
  "All 40 built-in scenarios",
  "8 languages",
  "Unlimited custom scenarios",
  "Instant grammar correction",
  "Vocabulary tracking",
  "All future updates",
];

export default function Pricing() {
  const productId = process.env.NEXT_PUBLIC_POLAR_PRODUCT_ID;
  const checkoutUrl = `/api/checkout?products=${productId}`;

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

        <div className="mx-auto max-w-lg">
          <div className="border border-border bg-bg-secondary p-6 sm:p-8">
            <div className="mb-4 text-sm text-yellow">
              A merchant approaches you...
            </div>

            <div className="mb-6 text-sm italic text-yellow">
              &ldquo;Greetings, adventurer. I have just what you need for your
              journey.&rdquo;
            </div>

            <div className="border-t border-border pt-6">
              <div className="mb-1 text-lg font-bold text-accent">
                LARK LICENSE
              </div>
              <div className="mb-6 flex items-baseline gap-1">
                <span className="text-3xl font-bold text-text">$2.99</span>
                <span className="text-sm text-text-dim">/month</span>
              </div>

              <div className="mb-8 space-y-2 text-sm">
                {features.map((feature) => (
                  <div key={feature} className="flex items-center gap-2">
                    <span className="text-accent">✓</span>
                    <span className="text-text">{feature}</span>
                  </div>
                ))}
              </div>

              <a
                href={checkoutUrl}
                className="block w-full border border-accent py-3 text-center text-sm font-bold text-accent transition-colors hover:bg-accent/10"
              >
                &gt; ACCEPT OFFER
              </a>
              <p className="mt-4 text-center text-xs text-text-dim">
                Cancel anytime &middot; Powered by Polar
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

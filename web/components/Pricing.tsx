"use client";

const features = [
  "All 25 built-in scenarios",
  "8 languages",
  "Unlimited custom scenarios",
  "Instant grammar correction",
  "Vocabulary tracking",
  "All future updates",
];

export default function Pricing() {
  const productId = process.env.NEXT_PUBLIC_POLAR_PRODUCT_ID;
  const checkoutUrl = productId
    ? `https://polar.sh/checkout?productId=${productId}`
    : "https://polar.sh";

  return (
    <section id="pricing" className="py-24">
      <div className="mx-auto max-w-6xl px-6">
        <h2 className="text-center text-3xl font-bold sm:text-4xl">
          Simple, affordable{" "}
          <span className="gradient-text">pricing</span>
        </h2>
        <p className="mx-auto mt-4 max-w-xl text-center text-text-dim">
          One plan. Everything included. Cancel anytime.
        </p>

        <div className="mx-auto mt-12 max-w-sm">
          <div className="rounded-lg border border-accent/30 bg-bg-card p-8">
            <div className="text-center">
              <div className="text-sm font-semibold uppercase tracking-wide text-accent">
                Monthly
              </div>
              <div className="mt-2 flex items-baseline justify-center gap-1">
                <span className="text-5xl font-bold">$3.99</span>
                <span className="text-text-dim">/month</span>
              </div>
            </div>

            <ul className="mt-8 space-y-3">
              {features.map((feature) => (
                <li key={feature} className="flex items-center gap-3 text-sm">
                  <span className="text-green">&#10003;</span>
                  {feature}
                </li>
              ))}
            </ul>

            <a
              href={checkoutUrl}
              className="mt-8 block w-full rounded-md bg-accent py-3 text-center font-semibold text-bg transition-colors hover:bg-accent-dim"
            >
              Subscribe Now
            </a>
            <p className="mt-3 text-center text-xs text-text-dim">
              Cancel anytime. Powered by Polar.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}

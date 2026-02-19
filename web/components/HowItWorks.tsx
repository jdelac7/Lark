const steps = [
  {
    number: "1",
    title: "Purchase",
    description: "Subscribe for $3.99/month. You'll receive a license key instantly after checkout.",
    code: null,
  },
  {
    number: "2",
    title: "Activate",
    description: "Run one command to activate your license and unlock Lark on your machine.",
    code: "$ lark activate YOUR-LICENSE-KEY",
  },
  {
    number: "3",
    title: "Play",
    description: "Launch Lark and start learning through immersive text adventures.",
    code: "$ lark",
  },
];

export default function HowItWorks() {
  return (
    <section id="how-it-works" className="border-y border-border bg-bg-secondary py-24">
      <div className="mx-auto max-w-6xl px-6">
        <h2 className="text-center text-3xl font-bold sm:text-4xl">
          Up and running in{" "}
          <span className="gradient-text">3 steps</span>
        </h2>

        <div className="mt-16 grid gap-8 lg:grid-cols-3">
          {steps.map((step) => (
            <div key={step.number} className="text-center">
              <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full border border-accent/30 bg-accent/10 text-lg font-bold text-accent">
                {step.number}
              </div>
              <h3 className="mt-4 text-xl font-semibold">{step.title}</h3>
              <p className="mt-2 text-sm text-text-dim">{step.description}</p>
              {step.code && (
                <code className="mt-4 inline-block rounded-md border border-border bg-bg px-4 py-2 text-left text-sm text-green">
                  {step.code}
                </code>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

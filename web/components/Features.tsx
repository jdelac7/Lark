const features = [
  {
    title: "25 Scenarios",
    description:
      "From restaurants to job interviews — practice real-world situations you'll actually encounter abroad.",
    icon: "🗺",
  },
  {
    title: "8 Languages",
    description:
      "Spanish, French, German, Japanese, Italian, Portuguese, Korean, and Chinese.",
    icon: "🌍",
  },
  {
    title: "Dynamic Stories",
    description:
      "Every conversation is unique. The game adapts to your choices and creates dynamic, branching storylines.",
    icon: "🎭",
  },
  {
    title: "Grammar Correction",
    description:
      "Write your own responses and get instant feedback on grammar, spelling, and natural phrasing.",
    icon: "✏",
  },
  {
    title: "Vocabulary Tracking",
    description:
      "Learn 2-4 new words per turn with translations and contextual usage notes.",
    icon: "📖",
  },
  {
    title: "Custom Scenarios",
    description:
      "Create your own scenarios — practice any situation you can imagine, from buying flowers to renting an apartment.",
    icon: "⚡",
  },
];

export default function Features() {
  return (
    <section id="features" className="py-24">
      <div className="mx-auto max-w-6xl px-6">
        <h2 className="text-center text-3xl font-bold sm:text-4xl">
          Everything you need to{" "}
          <span className="gradient-text">learn by doing</span>
        </h2>
        <p className="mx-auto mt-4 max-w-2xl text-center text-text-dim">
          Forget flashcards and grammar drills. Lark drops you into realistic
          scenarios where you learn by actually using the language.
        </p>

        <div className="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="rounded-lg border border-border bg-bg-card p-6 transition-colors hover:border-accent/30"
            >
              <div className="text-2xl">{feature.icon}</div>
              <h3 className="mt-4 text-lg font-semibold">{feature.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-text-dim">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

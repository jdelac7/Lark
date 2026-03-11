const features = [
  {
    label: "Scenarios",
    summary: "40 scenarios (25 everyday + 15 adventure)",
    detail:
      "From restaurants to job interviews — practice situations you'll actually encounter abroad.",
  },
  {
    label: "Languages",
    summary: "80+ supported",
    detail:
      "Spanish, French, German, Japanese, Arabic, Hindi, and 75 more — from Afrikaans to Zulu.",
  },
  {
    label: "Stories",
    summary: "Dynamic, branching narratives",
    detail:
      "Every playthrough is unique. The game adapts to your choices and creates branching storylines.",
  },
  {
    label: "Grammar",
    summary: "Instant correction & feedback",
    detail:
      "Write your own responses and get real-time feedback on grammar, spelling, and phrasing.",
  },
  {
    label: "Vocabulary",
    summary: "Contextual word tracking",
    detail:
      "Learn 2-4 new words per turn with translations and contextual usage notes.",
  },
  {
    label: "Custom",
    summary: "Create your own scenarios",
    detail:
      "Design any scenario you can imagine — from buying flowers to renting an apartment.",
  },
];

export default function Features() {
  return (
    <section id="features" className="border-t border-border py-20">
      <div className="mx-auto max-w-4xl px-6">
        {/* Command prompt */}
        <div className="mb-8 text-sm">
          <span className="text-accent">user@lark</span>
          <span className="text-text-dim">:</span>
          <span className="text-cyan">~</span>
          <span className="text-text-dim">$ </span>
          <span className="text-text">lark</span>
        </div>

        <div className="mb-8 text-sm text-text-dim">
          Forget flashcards and grammar drills. Lark drops you into realistic
          scenarios where you learn by actually using the language.
        </div>

        <div className="mb-4 text-sm font-bold text-yellow">FEATURES:</div>

        <div className="space-y-5 font-mono text-sm">
          {features.map((f) => (
            <div key={f.label}>
              <div className="flex flex-col gap-1 sm:flex-row sm:gap-0">
                <span className="w-44 shrink-0 text-accent">{f.label}</span>
                <span className="text-text">{f.summary}</span>
              </div>
              <div className="mt-1 text-xs leading-relaxed text-text-dim sm:pl-44">
                {f.detail}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

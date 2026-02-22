"use client";

import { useState, useEffect } from "react";

type Step =
  | { type: "narrative"; text: string; translation: string; delay: number }
  | { type: "npc"; text: string; translation: string; delay: number }
  | { type: "choices"; items: { text: string; translation: string }[]; delay: number }
  | { type: "input"; text: string; delay: number }
  | { type: "correction"; original: string; corrected: string; explanation: string; delay: number }
  | { type: "vocabulary"; words: { word: string; meaning: string }[]; delay: number };

const steps: Step[] = [
  {
    type: "narrative",
    text: "Entras en un acogedor restaurante de tapas en el barrio de La Latina de Madrid. El olor a ajo y aceite de oliva llena el aire. Un camarero se acerca con una cálida sonrisa.",
    translation: "You step into a cozy tapas restaurant in Madrid's La Latina neighborhood. The smell of garlic and olive oil fills the air. A waiter approaches with a warm smile.",
    delay: 200,
  },
  {
    type: "npc",
    text: "¡Buenas tardes! Bienvenido. ¿Mesa para cuántos?",
    translation: "Good afternoon! Welcome. Table for how many?",
    delay: 1200,
  },
  {
    type: "vocabulary",
    words: [
      { word: "buenas tardes", meaning: "good afternoon" },
      { word: "bienvenido", meaning: "welcome" },
      { word: "mesa", meaning: "table" },
    ],
    delay: 1800,
  },
  {
    type: "choices",
    items: [
      { text: "Mesa para uno, por favor", translation: "Table for one, please" },
      { text: "¿Tienen una terraza?", translation: "Do you have a terrace?" },
      { text: "¿Qué me recomienda?", translation: "What do you recommend?" },
    ],
    delay: 2200,
  },
  {
    type: "input",
    text: "Quiero una mesa para uno, por favour",
    delay: 3200,
  },
  {
    type: "correction",
    original: "por favour",
    corrected: "por favor",
    explanation: "\"Favor\" doesn't have a 'u' — a common English spelling influence.",
    delay: 4000,
  },
  {
    type: "narrative",
    text: "El camarero asiente y te lleva a una pequeña mesa junto a la ventana con vistas a la calle concurrida.",
    translation: "The waiter nods and leads you to a small table by the window overlooking the busy street below.",
    delay: 4800,
  },
  {
    type: "npc",
    text: "Aquí tiene la carta. ¿Le traigo algo de beber mientras mira?",
    translation: "Here's the menu. Can I bring you something to drink while you look?",
    delay: 5600,
  },
];

export default function TerminalMockup() {
  const [visibleSteps, setVisibleSteps] = useState(0);

  useEffect(() => {
    const timers: NodeJS.Timeout[] = [];
    steps.forEach((step, i) => {
      timers.push(
        setTimeout(() => setVisibleSteps(i + 1), step.delay)
      );
    });
    return () => timers.forEach(clearTimeout);
  }, []);

  const visible = steps.slice(0, visibleSteps);
  const animating = visibleSteps < steps.length;

  return (
    <div className="w-full overflow-hidden border border-accent/30 bg-bg-card">
      {/* Title bar */}
      <div className="flex items-center gap-2 border-b border-border px-4 py-2">
        <div className="h-2.5 w-2.5 rounded-full bg-accent/60" />
        <div className="h-2.5 w-2.5 rounded-full bg-yellow/60" />
        <div className="h-2.5 w-2.5 rounded-full bg-text-dim/40" />
        <span className="ml-2 text-xs text-text-dim">
          lark — At the Restaurant · Spanish
        </span>
      </div>

      {/* Content area */}
      <div className="p-4 font-mono text-sm">
        <div className="space-y-4">
          {visible.map((step, i) => {
            switch (step.type) {
              case "narrative":
                return (
                  <div key={i}>
                    <div className="whitespace-pre-wrap text-green">
                      {step.text}
                    </div>
                    <div className="mt-1 whitespace-pre-wrap text-text-dim italic">
                      {step.translation}
                    </div>
                  </div>
                );
              case "npc":
                return (
                  <div key={i} className="border-l-2 border-cyan/40 pl-3">
                    <div className="text-cyan">
                      &ldquo;{step.text}&rdquo;
                    </div>
                    <div className="mt-1 text-text-dim italic">
                      &ldquo;{step.translation}&rdquo;
                    </div>
                  </div>
                );
              case "vocabulary":
                return (
                  <div key={i} className="border-t border-border pt-3">
                    <div className="mb-2 text-xs font-bold text-purple">
                      VOCABULARY
                    </div>
                    <div className="space-y-1">
                      {step.words.map((w, j) => (
                        <div key={j} className="text-xs">
                          <span className="text-purple">{w.word}</span>
                          <span className="text-text-dim"> — </span>
                          <span className="text-text">{w.meaning}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              case "choices":
                return (
                  <div key={i} className="border-t border-border pt-3">
                    <div className="space-y-1">
                      {step.items.map((c, j) => (
                        <div key={j}>
                          <span className="text-text-dim">{j + 1}.</span>{" "}
                          <span className="text-cyan">{c.text}</span>
                          <span className="ml-2 text-text-dim italic">
                            ({c.translation})
                          </span>
                        </div>
                      ))}
                    </div>
                    <div className="mt-2 text-xs text-text-dim italic">
                      or write your own response...
                    </div>
                  </div>
                );
              case "input":
                return (
                  <div key={i} className="flex items-center border-t border-border pb-1 pt-3">
                    <span className="mr-2 text-accent">$</span>
                    <span className="text-text">{step.text}</span>
                  </div>
                );
              case "correction":
                return (
                  <div key={i} className="border border-yellow/30 bg-yellow/5 p-3">
                    <div className="mb-1 text-xs font-bold text-yellow">
                      CORRECTION
                    </div>
                    <div className="text-text-dim line-through">
                      {step.original}
                    </div>
                    <div className="text-green">{step.corrected}</div>
                    <div className="mt-1 text-xs text-yellow">
                      {step.explanation}
                    </div>
                  </div>
                );
              default:
                return null;
            }
          })}
        </div>

        {/* Blinking cursor while animating */}
        {animating && (
          <span className="cursor-blink mt-2 inline-block h-4 w-1.5 bg-accent" />
        )}
      </div>
    </div>
  );
}

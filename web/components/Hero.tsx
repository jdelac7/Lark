"use client";

import { useSession } from "next-auth/react";
import TerminalMockup from "./TerminalMockup";

export default function Hero() {
  const { data: session } = useSession();
  const isSubscribed = session?.user?.subscribed;

  return (
    <section className="relative min-h-screen pb-16 pt-20">
      <div className="mx-auto max-w-4xl px-6">
        {/* Command prompt */}
        <div className="mb-8 text-sm">
          <span className="text-accent">user@lark</span>
          <span className="text-text-dim">:</span>
          <span className="text-cyan">~</span>
          <span className="text-text-dim">$ </span>
          <span className="text-text">lark</span>
        </div>

        {/* ASCII Art Logo */}
        <pre className="text-glow-subtle text-sm leading-tight text-accent sm:text-base lg:text-lg">
          {[
            "  _                _    ",
            " | |    __ _  _ __| | __",
            " | |   / _` || '__| |/ /",
            " | |__| (_| || |  |   < ",
            " |_____\\__,_||_|  |_|\\_\\",
          ].join("\n")}
        </pre>

        <p className="mb-8 mt-5 text-sm text-text-dim">
          A Zork-inspired text-adventure game for learning languages.
        </p>

        <div className="mb-2 max-w-2xl text-sm leading-relaxed text-text">
          <p>
            You find yourself standing at the entrance of a language learning
            adventure. Before you lies a world of{" "}
            <span className="text-yellow">40 immersive scenarios</span> across{" "}
            <span className="text-cyan">80+ languages</span>. Practice ordering
            food in Madrid, navigating the streets of Tokyo, or checking into a
            hotel in Paris &mdash; all from your terminal.
          </p>
        </div>

        <p className="mb-10 text-sm text-text-dim">
          Type responses in your target language. Get instant grammar correction.
          Learn vocabulary in context. No flashcards. No drills. Just adventure.
        </p>

        {/* Game demo */}
        <div className="mb-10">
          <div className="mb-2 text-xs text-text-dim">
            <span className="text-accent">&gt;</span> demo
          </div>
          <TerminalMockup />
        </div>

        {/* CTAs as terminal commands */}
        <div className="flex flex-col gap-3 text-sm sm:flex-row">
          {isSubscribed ? (
            <a
              href="/play"
              onClick={() => {
                if (typeof window.gtag === "function") {
                  window.gtag("event", "cta_click", {
                    event_category: "engagement",
                    event_label: "hero_play",
                  });
                }
              }}
              className="inline-flex items-center gap-2 border border-accent px-5 py-2.5 text-accent transition-colors hover:bg-accent/10"
            >
              <span className="text-text-dim">&gt;</span> play
            </a>
          ) : (
            <a
              href="#shop"
              onClick={() => {
                if (typeof window.gtag === "function") {
                  window.gtag("event", "cta_click", {
                    event_category: "engagement",
                    event_label: "hero_subscribe",
                  });
                }
              }}
              className="inline-flex items-center gap-2 border border-accent px-5 py-2.5 text-accent transition-colors hover:bg-accent/10"
            >
              <span className="text-text-dim">&gt;</span> subscribe
            </a>
          )}
          <a
            href="#features"
            className="inline-flex items-center gap-2 border border-border px-5 py-2.5 text-text-dim transition-colors hover:border-text-dim hover:text-text"
          >
            <span className="text-text-dim">&gt;</span> help
          </a>
        </div>
      </div>
    </section>
  );
}

"use client";

import { useState, useEffect } from "react";

type Line =
  | { type: "plain"; text: string; color: string; delay: number }
  | { type: "header"; text: string; delay: number }
  | { type: "divider"; delay: number }
  | { type: "content"; text: string; color: string; delay: number }
  | { type: "empty"; delay: number };

const lines: Line[] = [
  { type: "plain", text: "$ lark", color: "text-green", delay: 0 },
  { type: "plain", text: "", color: "", delay: 300 },
  { type: "plain", text: "  _          _    ", color: "text-accent", delay: 400 },
  { type: "plain", text: " | |   __ _ | |__ ", color: "text-accent", delay: 500 },
  { type: "plain", text: " | |  / _` || '__|| |/ /", color: "text-accent", delay: 600 },
  { type: "plain", text: " | |_| (_| || |   |   < ", color: "text-accent", delay: 700 },
  { type: "plain", text: " |____\\__,_||_|   |_|\\_\\", color: "text-accent", delay: 800 },
  { type: "plain", text: "", color: "", delay: 900 },
  { type: "plain", text: "  A text-adventure language learning game", color: "text-text-dim", delay: 1000 },
  { type: "plain", text: "", color: "", delay: 1200 },
  // Box starts
  { type: "header", text: "Lark  ·  At the Restaurant  ·  Spanish", delay: 1400 },
  { type: "divider", delay: 1600 },
  { type: "empty", delay: 1700 },
  { type: "content", text: "You step into a cozy tapas restaurant", color: "text-[#60a5fa]", delay: 1800 },
  { type: "content", text: "in Madrid's La Latina neighborhood.", color: "text-[#60a5fa]", delay: 1900 },
  { type: "empty", delay: 2000 },
  { type: "content", text: 'NPC: "¡Buenas tardes! Bienvenido."', color: "text-yellow", delay: 2200 },
  { type: "content", text: '     "Good afternoon! Welcome."', color: "text-text-dim", delay: 2400 },
  { type: "empty", delay: 2500 },
  { type: "divider", delay: 2600 },
  { type: "content", text: '1) "Mesa para uno, por favor"', color: "text-green", delay: 2700 },
  { type: "content", text: '2) "¿Tienen una terraza?"', color: "text-green", delay: 2800 },
  { type: "content", text: "3) Write your own response...", color: "text-text-dim", delay: 2900 },
  // Box ends (handled by the container)
  { type: "plain", text: "> 1", color: "text-green", delay: 3400 },
];

export default function TerminalMockup() {
  const [visibleLines, setVisibleLines] = useState(0);

  useEffect(() => {
    const timers: NodeJS.Timeout[] = [];
    lines.forEach((line, i) => {
      timers.push(
        setTimeout(() => {
          setVisibleLines(i + 1);
        }, line.delay)
      );
    });
    return () => timers.forEach(clearTimeout);
  }, []);

  const visible = lines.slice(0, visibleLines);

  // Split into pre-box plain lines, box lines, and post-box plain lines
  const firstBoxIdx = visible.findIndex((l) => l.type !== "plain");
  const lastBoxIdx = visible.findLastIndex(
    (l) => l.type !== "plain"
  );

  const prePlain = firstBoxIdx === -1 ? visible : visible.slice(0, firstBoxIdx);
  const boxLines =
    firstBoxIdx === -1 ? [] : visible.slice(firstBoxIdx, lastBoxIdx + 1);
  const postPlain =
    firstBoxIdx === -1
      ? []
      : visible.slice(lastBoxIdx + 1).filter((l) => l.type === "plain");

  // Check if the box section has started
  const boxStarted = boxLines.length > 0;
  // Check if the last box line has appeared (the final content/divider before "> 1")
  const boxComplete =
    visibleLines >= lines.findLastIndex((l) => l.type !== "plain") + 1;

  return (
    <div className="w-full overflow-hidden rounded-lg border border-border bg-[#0d0d14] shadow-2xl shadow-accent/5">
      {/* Title bar */}
      <div className="flex items-center gap-2 border-b border-border px-4 py-3">
        <div className="h-3 w-3 rounded-full bg-[#ff5f57]" />
        <div className="h-3 w-3 rounded-full bg-[#ffbd2e]" />
        <div className="h-3 w-3 rounded-full bg-[#28c840]" />
        <span className="ml-2 text-xs text-text-dim">lark</span>
      </div>
      {/* Terminal content */}
      <div className="p-4 text-xs leading-5 sm:text-sm sm:leading-6">
        {/* Pre-box plain text (command + ASCII art) */}
        {prePlain.map((line, i) =>
          line.type === "plain" ? (
            <div key={i} className={`whitespace-pre font-mono ${line.color}`}>
              {line.text || "\u00A0"}
            </div>
          ) : null
        )}

        {/* Box rendered with CSS borders */}
        {boxStarted && (
          <div className="mt-1 border border-border/50 font-mono">
            {boxLines.map((line, i) => {
              switch (line.type) {
                case "header":
                  return (
                    <div
                      key={i}
                      className="border-b border-border/50 px-3 py-1 text-accent"
                    >
                      {line.text}
                    </div>
                  );
                case "divider":
                  return (
                    <div
                      key={i}
                      className="border-t border-border/50"
                    />
                  );
                case "empty":
                  return (
                    <div key={i} className="h-3" />
                  );
                case "content":
                  return (
                    <div
                      key={i}
                      className={`whitespace-pre px-3 ${line.color}`}
                    >
                      {line.text}
                    </div>
                  );
                default:
                  return null;
              }
            })}
          </div>
        )}

        {/* Post-box plain text ("> 1") */}
        {postPlain.map((line, i) =>
          line.type === "plain" ? (
            <div
              key={`post-${i}`}
              className={`whitespace-pre font-mono ${line.color}`}
            >
              {line.text || "\u00A0"}
            </div>
          ) : null
        )}

        {/* Blinking cursor while animating */}
        {visibleLines < lines.length && (
          <span className="cursor-blink mt-1 inline-block h-4 w-2 bg-accent" />
        )}
      </div>
    </div>
  );
}

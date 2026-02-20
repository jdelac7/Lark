"use client";

import { useState, useEffect } from "react";

type Line =
  | { type: "plain"; text: string; color: string; delay: number }
  | { type: "header"; text: string; delay: number }
  | { type: "divider"; delay: number }
  | { type: "content"; text: string; color: string; delay: number }
  | { type: "empty"; delay: number };

const lines: Line[] = [
  // Game box
  {
    type: "header",
    text: " At the Restaurant · Spanish · Madrid",
    delay: 200,
  },
  { type: "divider", delay: 400 },
  { type: "empty", delay: 500 },
  {
    type: "content",
    text: "  You step into a cozy tapas restaurant",
    color: "text-cyan",
    delay: 600,
  },
  {
    type: "content",
    text: "  in Madrid's La Latina neighborhood. The",
    color: "text-cyan",
    delay: 700,
  },
  {
    type: "content",
    text: "  smell of garlic and olive oil fills the air.",
    color: "text-cyan",
    delay: 800,
  },
  { type: "empty", delay: 900 },
  {
    type: "content",
    text: '  Camarero: "¡Buenas tardes! Bienvenido."',
    color: "text-yellow",
    delay: 1100,
  },
  {
    type: "content",
    text: '            "Good afternoon! Welcome."',
    color: "text-text-dim",
    delay: 1300,
  },
  { type: "empty", delay: 1400 },
  { type: "divider", delay: 1500 },
  {
    type: "content",
    text: '  1) "Mesa para uno, por favor"',
    color: "text-green",
    delay: 1600,
  },
  {
    type: "content",
    text: '  2) "¿Tienen una terraza?"',
    color: "text-green",
    delay: 1700,
  },
  {
    type: "content",
    text: "  3) Write your own response...",
    color: "text-text-dim",
    delay: 1800,
  },
  { type: "divider", delay: 1900 },
  // User input
  { type: "plain", text: "", color: "", delay: 2100 },
  { type: "plain", text: "> 1", color: "text-green", delay: 2400 },
  { type: "plain", text: "", color: "", delay: 2600 },
  // Response
  {
    type: "plain",
    text: '  You: "Mesa para uno, por favor"',
    color: "text-text",
    delay: 2800,
  },
  {
    type: "plain",
    text: "  ✓ Grammar: Correct!",
    color: "text-accent",
    delay: 3100,
  },
  { type: "plain", text: "", color: "", delay: 3200 },
  {
    type: "plain",
    text: "  + mesa (table)",
    color: "text-cyan",
    delay: 3400,
  },
  {
    type: "plain",
    text: "  + por favor (please)",
    color: "text-cyan",
    delay: 3600,
  },
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

  // Split into box lines and non-box lines
  const firstBoxIdx = visible.findIndex((l) => l.type !== "plain");
  const lastBoxIdx = visible.findLastIndex((l) => l.type !== "plain");

  const prePlain =
    firstBoxIdx === -1 ? visible : visible.slice(0, firstBoxIdx);
  const boxLines =
    firstBoxIdx === -1 ? [] : visible.slice(firstBoxIdx, lastBoxIdx + 1);
  const postPlain =
    firstBoxIdx === -1
      ? []
      : visible.slice(lastBoxIdx + 1).filter((l) => l.type === "plain");

  const boxStarted = boxLines.length > 0;

  return (
    <div className="w-full overflow-hidden border border-border bg-bg-secondary">
      {/* Title bar */}
      <div className="flex items-center border-b border-border px-4 py-2">
        <span className="text-xs text-text-dim">[lark · game demo]</span>
      </div>

      {/* Terminal content */}
      <div className="p-4 text-xs leading-5 sm:text-sm sm:leading-6">
        {/* Pre-box plain text */}
        {prePlain.map((line, i) =>
          line.type === "plain" ? (
            <div key={i} className={`whitespace-pre font-mono ${line.color}`}>
              {line.text || "\u00A0"}
            </div>
          ) : null
        )}

        {/* Game interface box */}
        {boxStarted && (
          <div className="border border-border/60 font-mono">
            {boxLines.map((line, i) => {
              switch (line.type) {
                case "header":
                  return (
                    <div
                      key={i}
                      className="border-b border-border/60 px-3 py-1.5 text-xs font-bold text-accent"
                    >
                      {line.text}
                    </div>
                  );
                case "divider":
                  return (
                    <div key={i} className="border-t border-border/60" />
                  );
                case "empty":
                  return <div key={i} className="h-3" />;
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

        {/* Post-box content (user input + response) */}
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
          <span className="cursor-blink mt-1 inline-block h-4 w-1.5 bg-accent" />
        )}
      </div>
    </div>
  );
}

"use client";

import { useState } from "react";

export default function CopyCmd({
  cmd,
  className = "",
}: {
  cmd: string;
  className?: string;
}) {
  const [copied, setCopied] = useState(false);

  function handleCopy() {
    navigator.clipboard.writeText(cmd).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
      if (typeof window.gtag === "function") {
        window.gtag("event", "copy_install_command", {
          event_category: "download",
          event_label: cmd,
        });
      }
    });
  }

  return (
    <button
      type="button"
      onClick={handleCopy}
      className={`group flex w-full cursor-pointer items-center gap-2 text-left ${className}`}
      title="Copy to clipboard"
    >
      <span className="min-w-0 flex-1 truncate">$ {cmd}</span>
      <span className="shrink-0 text-text-dim opacity-0 transition-opacity group-hover:opacity-100">
        {copied ? (
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-accent">
            <polyline points="20 6 9 17 4 12" />
          </svg>
        ) : (
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
          </svg>
        )}
      </span>
    </button>
  );
}

"use client";

import { useState } from "react";

export default function FeedbackWidget() {
  const [open, setOpen] = useState(false);
  const [message, setMessage] = useState("");
  const [status, setStatus] = useState<"idle" | "sending" | "sent" | "error">(
    "idle"
  );

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!message.trim()) return;

    setStatus("sending");
    try {
      const res = await fetch("/api/feedback", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: message.trim() }),
      });
      if (!res.ok) throw new Error();
      setStatus("sent");
      setMessage("");
      setTimeout(() => {
        setOpen(false);
        setStatus("idle");
      }, 2000);
    } catch {
      setStatus("error");
    }
  }

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="fixed bottom-5 right-5 z-50 cursor-pointer rounded border border-accent/50 bg-accent/10 px-3 py-1.5 font-mono text-xs text-accent transition-colors hover:border-accent hover:bg-accent/20"
      >
        [feedback]
      </button>

      {open && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
          <div className="w-full max-w-md rounded border border-border bg-bg p-6">
            <div className="mb-4 text-sm">
              <span className="text-accent">user@lark</span>
              <span className="text-text-dim">:</span>
              <span className="text-cyan">~</span>
              <span className="text-text-dim">$ </span>
              <span className="text-text">send-feedback</span>
            </div>

            {status === "sent" ? (
              <p className="text-sm text-green">
                feedback received — thanks!
              </p>
            ) : (
              <form onSubmit={handleSubmit}>
                <textarea
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder="bug report, feature request, or general feedback..."
                  maxLength={2000}
                  rows={5}
                  className="mb-3 w-full resize-none rounded border border-border bg-bg-secondary px-3 py-2 font-mono text-sm text-text placeholder:text-text-dim/50 focus:border-accent focus:outline-none"
                  autoFocus
                />
                <div className="flex items-center justify-between">
                  <span className="text-xs text-text-dim">
                    {message.length}/2000
                  </span>
                  <div className="flex gap-2">
                    <button
                      type="button"
                      onClick={() => {
                        setOpen(false);
                        setStatus("idle");
                      }}
                      className="cursor-pointer rounded border border-border px-3 py-1.5 text-xs text-text-dim transition-colors hover:text-text"
                    >
                      cancel
                    </button>
                    <button
                      type="submit"
                      disabled={
                        !message.trim() || status === "sending"
                      }
                      className="cursor-pointer rounded border border-accent bg-accent/10 px-3 py-1.5 text-xs text-accent transition-colors hover:bg-accent/20 disabled:opacity-40"
                    >
                      {status === "sending" ? "sending..." : "send"}
                    </button>
                  </div>
                </div>
                {status === "error" && (
                  <p className="mt-2 text-xs text-red-400">
                    failed to send — please try again
                  </p>
                )}
              </form>
            )}
          </div>
        </div>
      )}
    </>
  );
}

"use client";

import { signIn } from "next-auth/react";
import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

interface AuthFormProps {
  mode: "login" | "register";
}

export default function AuthForm({ mode }: AuthFormProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const callbackUrl = searchParams.get("callbackUrl") || "/play";
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const isRegister = mode === "register";

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      if (isRegister) {
        const res = await fetch("/api/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password, name }),
        });

        if (!res.ok) {
          const data = await res.json();
          setError(data.error || "Registration failed");
          setLoading(false);
          return;
        }
      }

      const result = await signIn("credentials", {
        email,
        password,
        redirect: false,
      });

      if (result?.error) {
        setError("Invalid email or password");
        setLoading(false);
        return;
      }

      router.push(callbackUrl);
      router.refresh();
    } catch {
      setError("Something went wrong. Try again.");
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-6">
      <div className="w-full max-w-md">
        <div className="border border-accent/30 bg-bg-card p-8">
          {/* Header */}
          <div className="mb-6">
            <div className="mb-2 text-sm text-text-dim">
              <span className="text-accent">user@lark</span>
              <span className="text-text-dim">:</span>
              <span className="text-cyan">~</span>
              <span className="text-text-dim">$ </span>
              <span className="text-text">
                lark {isRegister ? "register" : "login"}
              </span>
            </div>
            <h1 className="text-xl font-bold text-accent">
              {isRegister ? "CREATE ACCOUNT" : "LOGIN"}
            </h1>
          </div>

          {error && (
            <div className="mb-4 border border-yellow/30 bg-yellow/10 px-3 py-2 text-sm text-yellow">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            {isRegister && (
              <div>
                <label className="mb-1 block text-xs text-text-dim">
                  name:
                </label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full border border-border bg-bg px-3 py-2 font-mono text-sm text-text outline-none focus:border-accent"
                  placeholder="adventurer"
                />
              </div>
            )}

            <div>
              <label className="mb-1 block text-xs text-text-dim">
                email:
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="w-full border border-border bg-bg px-3 py-2 font-mono text-sm text-text outline-none focus:border-accent"
                placeholder="you@example.com"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs text-text-dim">
                password:
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={6}
                className="w-full border border-border bg-bg px-3 py-2 font-mono text-sm text-text outline-none focus:border-accent"
                placeholder="••••••••"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full border border-accent py-3 text-sm font-bold text-accent transition-colors hover:bg-accent/10 disabled:opacity-50"
            >
              {loading
                ? "> PROCESSING..."
                : isRegister
                  ? "> CREATE ACCOUNT"
                  : "> LOGIN"}
            </button>
          </form>

          <div className="my-6 flex items-center gap-3">
            <div className="h-px flex-1 bg-border" />
            <span className="text-xs text-text-dim">or</span>
            <div className="h-px flex-1 bg-border" />
          </div>

          <button
            type="button"
            onClick={() => signIn("google", { callbackUrl })}
            className="flex w-full items-center justify-center gap-2 border border-border py-3 text-sm text-text transition-colors hover:border-accent hover:text-accent"
          >
            <svg className="h-4 w-4" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.27-4.74 3.27-8.1z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Continue with Google
          </button>

          <div className="mt-6 text-center text-xs text-text-dim">
            {isRegister ? (
              <>
                Already have an account?{" "}
                <a href={`/login${callbackUrl !== "/play" ? `?callbackUrl=${encodeURIComponent(callbackUrl)}` : ""}`} className="text-accent hover:text-accent-dim">
                  [login]
                </a>
              </>
            ) : (
              <>
                Need an account?{" "}
                <a
                  href={`/register${callbackUrl !== "/play" ? `?callbackUrl=${encodeURIComponent(callbackUrl)}` : ""}`}
                  className="text-accent hover:text-accent-dim"
                >
                  [register]
                </a>
              </>
            )}
          </div>

          <div className="mt-4 text-center">
            <a
              href="/"
              className="text-xs text-text-dim transition-colors hover:text-accent"
            >
              &larr; back to home
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

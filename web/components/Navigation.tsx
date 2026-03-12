"use client";

import { useSession, signOut } from "next-auth/react";

export default function Navigation() {
  const { data: session } = useSession();

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-border bg-bg/90 backdrop-blur-sm">
      <div className="mx-auto flex max-w-4xl items-center justify-between px-6 py-3">
        <div className="flex items-center gap-3">
          <a href="/" className="text-glow-subtle font-bold text-accent">
            lark
          </a>
          <span className="text-xs text-text-dim">v1.0.2</span>
          <span className="border border-yellow/40 bg-yellow/10 px-1.5 py-0.5 text-[10px] font-bold uppercase text-yellow">
            beta
          </span>
        </div>
        <div className="flex items-center text-xs text-text-dim">
          <a
            href="#features"
            className="px-2.5 py-1 transition-colors hover:text-accent"
          >
            [features]
          </a>
          <a
            href="#quickstart"
            className="hidden px-2.5 py-1 transition-colors hover:text-accent sm:inline"
          >
            [quickstart]
          </a>
          {!session?.user?.subscribed && (
            <a
              href="#shop"
              className="px-2.5 py-1 transition-colors hover:text-accent"
            >
              [shop]
            </a>
          )}
          {session?.user?.subscribed ? (
            <a
              href="/play"
              className="px-2.5 py-1 font-bold text-accent transition-colors hover:text-accent-dim"
            >
              [play]
            </a>
          ) : session ? (
            <a
              href="/#shop"
              className="px-2.5 py-1 text-yellow transition-colors hover:text-accent"
            >
              [subscribe]
            </a>
          ) : (
            <a
              href="/login"
              className="px-2.5 py-1 transition-colors hover:text-accent"
            >
              [login]
            </a>
          )}
          {session?.user?.isAdmin && (
            <a
              href="/admin"
              className="px-2.5 py-1 text-cyan transition-colors hover:text-accent"
            >
              [admin]
            </a>
          )}
          {session && (
            <button
              onClick={() => signOut({ callbackUrl: "/" })}
              className="cursor-pointer px-2.5 py-1 text-text-dim transition-colors hover:text-accent"
            >
              [logout]
            </button>
          )}
        </div>
      </div>
    </nav>
  );
}

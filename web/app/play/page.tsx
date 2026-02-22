import { auth } from "@/lib/auth";
import GameTerminal from "./GameTerminal";

export const metadata = {
  title: "Play - Lark",
};

export default async function PlayPage() {
  const session = await auth();

  if (!session?.user?.subscribed) {
    return (
      <div className="flex min-h-screen items-center justify-center px-6">
        <div className="w-full max-w-md">
          <div className="border border-accent/30 bg-bg-card p-8">
            <div className="mb-2 text-sm text-text-dim">
              <span className="text-accent">user@lark</span>
              <span className="text-text-dim">:</span>
              <span className="text-cyan">~</span>
              <span className="text-text-dim">$ </span>
              <span className="text-text">lark play</span>
            </div>
            <div className="mb-4 text-lg font-bold text-yellow">
              ACCESS DENIED
            </div>
            <p className="mb-6 text-sm text-text-dim">
              An active subscription is required to play in the browser.
            </p>
            <div className="space-y-2">
              <a
                href="/#shop"
                className="block w-full border border-accent py-3 text-center text-sm font-bold text-accent transition-colors hover:bg-accent/10"
              >
                &gt; SUBSCRIBE
              </a>
              {!session && (
                <a
                  href="/login"
                  className="block w-full border border-cyan py-3 text-center text-sm font-bold text-cyan transition-colors hover:bg-cyan/10"
                >
                  &gt; LOGIN
                </a>
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

  return <GameTerminal />;
}

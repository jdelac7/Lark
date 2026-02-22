import { redirect } from "next/navigation";
import { auth } from "@/lib/auth";

export default async function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = await auth();

  if (!session?.user?.isAdmin) {
    redirect("/");
  }

  return (
    <div className="min-h-screen bg-bg">
      <header className="border-b border-border bg-bg-secondary/80 backdrop-blur-sm">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-3">
          <div className="flex items-center gap-3">
            <a
              href="/"
              className="text-glow-subtle font-bold text-accent"
            >
              lark
            </a>
            <span className="text-text-dim">/</span>
            <span className="text-sm text-cyan">admin</span>
          </div>
          <div className="flex items-center gap-4 text-xs text-text-dim">
            <span>{session.user.email}</span>
            <a
              href="/"
              className="transition-colors hover:text-accent"
            >
              [exit]
            </a>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-6 py-8">{children}</main>
    </div>
  );
}

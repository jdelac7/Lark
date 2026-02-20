export default function Navigation() {
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-border bg-bg/90 backdrop-blur-sm">
      <div className="mx-auto flex max-w-4xl items-center justify-between px-6 py-3">
        <div className="flex items-center gap-3">
          <span className="text-glow-subtle font-bold text-accent">lark</span>
          <span className="text-xs text-text-dim">v0.8.0</span>
          <span className="border border-yellow/40 bg-yellow/10 px-1.5 py-0.5 text-[10px] font-bold uppercase text-yellow">beta</span>
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
          <a
            href="#shop"
            className="px-2.5 py-1 transition-colors hover:text-accent"
          >
            [shop]
          </a>
        </div>
      </div>
    </nav>
  );
}

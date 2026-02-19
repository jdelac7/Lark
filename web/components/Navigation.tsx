export default function Navigation() {
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-border bg-bg/80 backdrop-blur-md">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <a href="#" className="text-lg font-bold text-accent">
          Lark
        </a>
        <div className="flex items-center gap-6 text-sm text-text-dim">
          <a
            href="#features"
            className="hidden transition-colors hover:text-text sm:inline"
          >
            Features
          </a>
          <a
            href="#how-it-works"
            className="hidden transition-colors hover:text-text sm:inline"
          >
            How It Works
          </a>
          <a
            href="#pricing"
            className="rounded-md bg-accent px-4 py-2 font-semibold text-bg transition-colors hover:bg-accent-dim"
          >
            Get Started
          </a>
        </div>
      </div>
    </nav>
  );
}

export default function Footer() {
  return (
    <footer className="border-t border-border py-12">
      <div className="mx-auto max-w-4xl px-6">
        {/* Command prompt */}
        <div className="mb-4 text-sm">
          <span className="text-accent">user@lark</span>
          <span className="text-text-dim">:</span>
          <span className="text-cyan">~</span>
          <span className="text-text-dim">$ </span>
          <span className="text-text">exit</span>
        </div>

        <p className="mb-6 text-sm text-text-dim">
          Thanks for visiting. May your adventures be fruitful.
        </p>

        <div className="flex flex-col justify-between gap-4 text-xs text-text-dim sm:flex-row sm:items-center">
          <div suppressHydrationWarning>
            <span className="text-accent">lark</span> &copy;{" "}
            {new Date().getFullYear()}
          </div>
          <div className="flex gap-4">
<a
              href="#shop"
              className="transition-colors hover:text-text"
            >
              [pricing]
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}

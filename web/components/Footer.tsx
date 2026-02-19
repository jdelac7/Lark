export default function Footer() {
  return (
    <footer className="border-t border-border py-8">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 sm:flex-row">
        <div className="text-sm text-text-dim">
          <span className="font-semibold text-accent">Lark</span>{" "}
          &copy; {new Date().getFullYear()}
        </div>
        <div className="flex gap-6 text-sm text-text-dim">
          <a
            href="https://github.com/joshburnsxyz/lark"
            className="transition-colors hover:text-text"
            target="_blank"
            rel="noopener noreferrer"
          >
            GitHub
          </a>
          <a href="#pricing" className="transition-colors hover:text-text">
            Pricing
          </a>
        </div>
      </div>
    </footer>
  );
}

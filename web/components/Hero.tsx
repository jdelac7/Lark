import TerminalMockup from "./TerminalMockup";

export default function Hero() {
  return (
    <section className="relative flex min-h-screen items-center pt-20">
      {/* Background gradient */}
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--color-accent)_0%,_transparent_50%)] opacity-[0.07]" />

      <div className="mx-auto grid max-w-6xl gap-12 px-6 lg:grid-cols-2 lg:gap-16">
        {/* Left: Copy */}
        <div className="flex flex-col justify-center">
          <h1 className="text-4xl font-bold leading-tight sm:text-5xl lg:text-6xl">
            Learn Languages
            <br />
            <span className="gradient-text">Through Adventure</span>
          </h1>
          <p className="mt-6 max-w-lg text-lg text-text-dim">
            A text-adventure game that teaches you real-world language skills
            through immersive role-play scenarios. Practice ordering food,
            checking into hotels, navigating cities, and more — all in your
            target language.
          </p>
          <div className="mt-8 flex gap-4">
            <a
              href="#pricing"
              className="rounded-md bg-accent px-6 py-3 font-semibold text-bg transition-colors hover:bg-accent-dim"
            >
              Get Started
            </a>
            <a
              href="#features"
              className="rounded-md border border-border px-6 py-3 font-semibold text-text-dim transition-colors hover:border-text-dim hover:text-text"
            >
              Learn More
            </a>
          </div>
        </div>

        {/* Right: Terminal */}
        <div className="flex items-center">
          <TerminalMockup />
        </div>
      </div>
    </section>
  );
}

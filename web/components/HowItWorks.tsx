export default function HowItWorks() {
  return (
    <section id="quickstart" className="border-t border-border py-20">
      <div className="mx-auto max-w-4xl px-6">
        {/* Command prompt */}
        <div className="mb-8 text-sm">
          <span className="text-accent">user@lark</span>
          <span className="text-text-dim">:</span>
          <span className="text-cyan">~</span>
          <span className="text-text-dim">$ </span>
          <span className="text-text">lark quickstart</span>
        </div>

        <div className="mb-8 text-sm font-bold text-yellow">
          QUICKSTART GUIDE
        </div>

        <div className="space-y-8 text-sm">
          {/* Step 1 */}
          <div>
            <div className="mb-2 text-text-dim">
              <span className="text-accent">Step 1</span> &mdash; Install
            </div>
            <div className="border border-border bg-bg-secondary p-4">
              <div className="text-green">
                $ curl -fsSL https://lark.joshburns.xyz/install.sh | sh
              </div>
              <div className="mt-2 text-accent">
                ✓ Lark installed to ~/.local/bin/lark
              </div>
            </div>
          </div>

          {/* Step 2 */}
          <div>
            <div className="mb-2 text-text-dim">
              <span className="text-accent">Step 2</span> &mdash; Subscribe
            </div>
            <div className="border border-border bg-bg-secondary p-4">
              <div className="text-green">$ open https://lark.joshburns.xyz</div>
              <div className="mt-2 text-text-dim">
                &rarr; Subscribe for $2.99/month
              </div>
              <div className="text-text-dim">
                &rarr; Receive your license key instantly
              </div>
            </div>
          </div>

          {/* Step 3 */}
          <div>
            <div className="mb-2 text-text-dim">
              <span className="text-accent">Step 3</span> &mdash; Activate
              &amp; Play
            </div>
            <div className="border border-border bg-bg-secondary p-4">
              <div className="text-green">
                $ lark activate YOUR-LICENSE-KEY
              </div>
              <div className="mt-2 text-accent">
                ✓ License activated successfully
              </div>
              <div className="mt-2 text-green">$ lark</div>
              <div className="mt-2 text-accent">
                &nbsp;&nbsp;Starting Lark v0.8.0...
              </div>
              <div className="text-text-dim">
                &nbsp;&nbsp;Loading scenarios...
              </div>
              <div className="mt-1 text-cyan">
                &nbsp;&nbsp;Welcome back, adventurer. Where shall we go today?
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

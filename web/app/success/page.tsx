import { polar } from "@/lib/polar";

interface SuccessPageProps {
  searchParams: Promise<{ checkout_id?: string }>;
}

export default async function SuccessPage({ searchParams }: SuccessPageProps) {
  const { checkout_id } = await searchParams;

  let licenseKey: string | null = null;
  let error: string | null = null;

  if (checkout_id) {
    try {
      // Fetch checkout to get customer ID
      const checkout = await polar.checkouts.get({ id: checkout_id });

      if (checkout.customerId) {
        // Use server-side API to list all license keys, then filter by customer
        const orgId = process.env.POLAR_ORGANIZATION_ID;
        const keys = await polar.licenseKeys.list({
          organizationId: orgId,
          limit: 100,
        });

        const customerKey = keys.result.items.find(
          (k) => k.customerId === checkout.customerId
        );
        if (customerKey) {
          licenseKey = customerKey.key;
        }
      }
    } catch (e) {
      error =
        "Unable to retrieve your license key. Please check your email for the key, or contact support.";
      console.error("Failed to fetch checkout:", e);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-6">
      <div className="w-full max-w-lg">
        <div className="rounded-lg border border-accent/30 bg-bg-card p-8">
          <div className="text-center">
            <div className="text-4xl text-green">&#10003;</div>
            <h1 className="mt-4 text-2xl font-bold">Welcome to Lark!</h1>
            <p className="mt-2 text-text-dim">
              Your subscription is active. Here&apos;s how to get started:
            </p>
          </div>

          {error && (
            <div className="mt-6 rounded-md border border-yellow/30 bg-yellow/10 p-4 text-sm text-yellow">
              {error}
            </div>
          )}

          {licenseKey && (
            <div className="mt-6">
              <label className="text-sm font-semibold text-text-dim">
                Your License Key
              </label>
              <div className="mt-2 flex items-center gap-2 rounded-md border border-border bg-bg p-3">
                <code className="flex-1 break-all text-sm text-accent">
                  {licenseKey}
                </code>
              </div>
            </div>
          )}

          <div className="mt-8 space-y-4">
            <h2 className="font-semibold">Activation Instructions</h2>
            <div className="space-y-3 text-sm">
              <div className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-bold text-accent">
                  1
                </span>
                <div>
                  <p>Install Lark (if you haven&apos;t already):</p>
                  <code className="mt-1 block rounded border border-border bg-bg px-3 py-2 text-green">
                    go install github.com/joshburnsxyz/lark/cli@latest
                  </code>
                </div>
              </div>
              <div className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-bold text-accent">
                  2
                </span>
                <div>
                  <p>Activate your license:</p>
                  <code className="mt-1 block rounded border border-border bg-bg px-3 py-2 text-green">
                    lark activate{" "}
                    {licenseKey ? licenseKey : "<YOUR-LICENSE-KEY>"}
                  </code>
                </div>
              </div>
              <div className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/10 text-xs font-bold text-accent">
                  3
                </span>
                <div>
                  <p>Start playing:</p>
                  <code className="mt-1 block rounded border border-border bg-bg px-3 py-2 text-green">
                    lark
                  </code>
                </div>
              </div>
            </div>
          </div>

          <div className="mt-8 text-center">
            <a
              href="/"
              className="text-sm text-accent transition-colors hover:text-accent-dim"
            >
              &larr; Back to home
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}

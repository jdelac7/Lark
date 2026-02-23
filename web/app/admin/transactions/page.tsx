import { getAllUsers } from "@/lib/db";
import { isCostDbAvailable, getAllCostEvents } from "@/lib/cost-db";
import TransactionsTable from "../components/TransactionsTable";
import Link from "next/link";

export const dynamic = "force-dynamic";

export default function TransactionsPage() {
  const costAvailable = isCostDbAvailable();
  const transactions = costAvailable ? getAllCostEvents() : [];
  const users = getAllUsers();

  const emailMap: Record<string, string> = {};
  for (const u of users) {
    emailMap[u.id] = u.email;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-lg font-bold text-accent">all transactions</h1>
          <p className="text-xs text-text-dim">
            every openrouter api call and its cost
          </p>
        </div>
        <Link
          href="/admin"
          className="rounded border border-border px-3 py-1.5 text-xs text-text-dim transition-colors hover:bg-bg-secondary hover:text-text"
        >
          back to dashboard
        </Link>
      </div>

      {!costAvailable && (
        <div className="rounded border border-yellow/30 bg-yellow/5 px-4 py-3 text-sm text-yellow">
          cost.db not found — transaction data is unavailable.
        </div>
      )}

      <TransactionsTable transactions={transactions} emailMap={emailMap} />
    </div>
  );
}

"use client";

import { useState, useMemo } from "react";

interface Transaction {
  id: number;
  player_id: string;
  scenario_id: string;
  language: string;
  is_custom: number;
  amount: number;
  created_at: string;
}

const PAGE_SIZE = 50;

export default function TransactionsTable({
  transactions,
  emailMap,
}: {
  transactions: Transaction[];
  emailMap: Record<string, string>;
}) {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    if (!search) return transactions;
    const q = search.toLowerCase();
    return transactions.filter(
      (t) =>
        t.player_id.toLowerCase().includes(q) ||
        (emailMap[t.player_id] ?? "").toLowerCase().includes(q) ||
        t.scenario_id.toLowerCase().includes(q) ||
        t.language.toLowerCase().includes(q)
    );
  }, [transactions, search, emailMap]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE));
  const clamped = Math.min(page, totalPages - 1);
  const slice = filtered.slice(clamped * PAGE_SIZE, (clamped + 1) * PAGE_SIZE);
  const totalCost = filtered.reduce((sum, t) => sum + t.amount, 0);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between gap-4">
        <input
          type="text"
          placeholder="search by player, email, scenario, language..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(0);
          }}
          className="w-full max-w-sm rounded border border-border bg-bg-secondary px-3 py-1.5 text-xs text-text placeholder:text-text-dim focus:border-accent focus:outline-none"
        />
        <div className="shrink-0 text-xs text-text-dim">
          {filtered.length} transactions &middot; total{" "}
          <span className="font-mono text-yellow">
            ${totalCost.toFixed(4)}
          </span>
        </div>
      </div>

      <div className="overflow-x-auto rounded border border-border">
        <table className="w-full text-left text-xs">
          <thead className="border-b border-border bg-bg-secondary text-text-dim">
            <tr>
              <th className="px-4 py-2">id</th>
              <th className="px-4 py-2">player</th>
              <th className="px-4 py-2">scenario</th>
              <th className="px-4 py-2">language</th>
              <th className="px-4 py-2">type</th>
              <th className="px-4 py-2 text-right">cost</th>
              <th className="px-4 py-2">time</th>
            </tr>
          </thead>
          <tbody>
            {slice.map((t) => (
              <tr
                key={t.id}
                className="border-b border-border/50 transition-colors hover:bg-bg-secondary/50"
              >
                <td className="px-4 py-2 font-mono text-text-dim">{t.id}</td>
                <td className="px-4 py-2 text-text">
                  {emailMap[t.player_id] ?? t.player_id.slice(0, 12) + "..."}
                </td>
                <td className="px-4 py-2 text-text">{t.scenario_id}</td>
                <td className="px-4 py-2 text-cyan">
                  {t.language || "---"}
                </td>
                <td className="px-4 py-2">
                  {t.is_custom ? (
                    <span className="text-purple">custom</span>
                  ) : (
                    <span className="text-green">premade</span>
                  )}
                </td>
                <td className="px-4 py-2 text-right font-mono text-yellow">
                  ${t.amount.toFixed(6)}
                </td>
                <td className="px-4 py-2 text-text-dim">
                  {t.created_at?.replace("T", " ").slice(0, 19) ?? "---"}
                </td>
              </tr>
            ))}
            {slice.length === 0 && (
              <tr>
                <td
                  colSpan={7}
                  className="px-4 py-8 text-center text-text-dim"
                >
                  {search ? "no matching transactions" : "no transactions yet"}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between text-xs">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={clamped === 0}
            className="rounded border border-border px-3 py-1 text-text-dim transition-colors hover:bg-bg-secondary disabled:opacity-30"
          >
            prev
          </button>
          <span className="text-text-dim">
            page {clamped + 1} of {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
            disabled={clamped >= totalPages - 1}
            className="rounded border border-border px-3 py-1 text-text-dim transition-colors hover:bg-bg-secondary disabled:opacity-30"
          >
            next
          </button>
        </div>
      )}
    </div>
  );
}

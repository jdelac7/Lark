import { getAllUsers, getFeedbackCount } from "@/lib/db";
import {
  isCostDbAvailable,
  getTotalCostAllPlayers,
  getPlayerCosts,
  getCostEventCount,
  getDailyCostSplit,
  getTopScenariosByCost,
  getUniquePlayerCount,
  getCostSplitTotals,
} from "@/lib/cost-db";
import Link from "next/link";
import AdminStats from "./components/AdminStats";
import UsersTable from "./components/UsersTable";
import CostChart from "./components/CostChart";

export const dynamic = "force-dynamic";

export default function AdminPage() {
  const costAvailable = isCostDbAvailable();
  const users = getAllUsers();
  const feedbackCount = getFeedbackCount();

  const totalCost = costAvailable ? getTotalCostAllPlayers() : 0;
  const playerCosts = costAvailable ? getPlayerCosts() : [];
  const eventCount = costAvailable ? getCostEventCount() : 0;
  const dailyCosts = costAvailable ? getDailyCostSplit(30) : [];
  const topScenarios = costAvailable ? getTopScenariosByCost(20) : [];
  const activePlayers = costAvailable ? getUniquePlayerCount() : 0;
  const costSplit = costAvailable
    ? getCostSplitTotals()
    : { premade: 0, custom: 0 };

  const subscribedCount = users.filter((u) => u.subscribed === 1).length;
  const subRate =
    users.length > 0
      ? ((subscribedCount / users.length) * 100).toFixed(1)
      : "0";

  // Build cost lookup by player_id
  const costMap = new Map(playerCosts.map((p) => [p.player_id, p.total_cost]));

  const userRows = users.map((u) => ({
    id: u.id,
    email: u.email,
    name: u.name,
    subscribed: u.subscribed,
    is_admin: u.is_admin,
    created_at: u.created_at,
    ai_cost: costMap.get(u.id) ?? 0,
  }));

  const stats = [
    { label: "total users", value: users.length, color: "text-accent" },
    {
      label: "subscribed",
      value: subscribedCount,
      sub: `${subRate}% rate`,
      color: "text-green",
    },
    {
      label: "total ai cost",
      value: `$${totalCost.toFixed(4)}`,
      color: "text-yellow",
    },
    { label: "api requests", value: eventCount, color: "text-cyan" },
    { label: "active players", value: activePlayers, color: "text-purple" },
    {
      label: "premade cost",
      value: `$${costSplit.premade.toFixed(4)}`,
      color: "text-green",
    },
    {
      label: "custom cost",
      value: `$${costSplit.custom.toFixed(4)}`,
      color: "text-purple",
    },
    {
      label: "feedback",
      value: feedbackCount,
      color: "text-yellow",
    },
  ];

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-lg font-bold text-accent">dashboard</h1>
          <p className="text-xs text-text-dim">system overview and analytics</p>
        </div>
        <div className="flex gap-2">
          <Link
            href="/admin/feedback"
            className="rounded border border-border bg-bg-secondary px-4 py-2 text-xs font-medium text-text transition-colors hover:border-accent hover:text-accent"
          >
            feedback
          </Link>
          <Link
            href="/admin/transactions"
            className="rounded border border-border bg-bg-secondary px-4 py-2 text-xs font-medium text-text transition-colors hover:border-accent hover:text-accent"
          >
            all transactions
          </Link>
        </div>
      </div>

      {!costAvailable && (
        <div className="rounded border border-yellow/30 bg-yellow/5 px-4 py-3 text-sm text-yellow">
          cost.db not found — cost analytics are unavailable. Set COST_DB_PATH
          or ensure the Go server has been run at least once.
        </div>
      )}

      <AdminStats stats={stats} />

      <section>
        <h2 className="mb-3 text-sm font-bold text-text">
          cost over time{" "}
          <span className="font-normal text-text-dim">(30 days)</span>
        </h2>
        <CostChart data={dailyCosts} />
      </section>

      <section>
        <h2 className="mb-3 text-sm font-bold text-text">
          top scenarios{" "}
          <span className="font-normal text-text-dim">by cost</span>
        </h2>
        <div className="overflow-x-auto rounded border border-border">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-border bg-bg-secondary text-text-dim">
              <tr>
                <th className="px-4 py-2">scenario</th>
                <th className="px-4 py-2">language</th>
                <th className="px-4 py-2">type</th>
                <th className="px-4 py-2 text-right">requests</th>
                <th className="px-4 py-2 text-right">total cost</th>
              </tr>
            </thead>
            <tbody>
              {topScenarios.map((s, i) => (
                <tr
                  key={`${s.scenario_id}-${s.language}-${s.is_custom}-${i}`}
                  className="border-b border-border/50 transition-colors hover:bg-bg-secondary/50"
                >
                  <td className="px-4 py-2 text-text">{s.scenario_id}</td>
                  <td className="px-4 py-2 text-cyan">{s.language || "---"}</td>
                  <td className="px-4 py-2">
                    {s.is_custom ? (
                      <span className="text-purple">custom</span>
                    ) : (
                      <span className="text-green">premade</span>
                    )}
                  </td>
                  <td className="px-4 py-2 text-right text-text-dim">
                    {s.event_count}
                  </td>
                  <td className="px-4 py-2 text-right font-mono text-yellow">
                    ${s.total_cost.toFixed(4)}
                  </td>
                </tr>
              ))}
              {topScenarios.length === 0 && (
                <tr>
                  <td
                    colSpan={5}
                    className="px-4 py-8 text-center text-text-dim"
                  >
                    no scenario data yet
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="mb-3 text-sm font-bold text-text">users</h2>
        <UsersTable users={userRows} />
      </section>
    </div>
  );
}

import Database from "better-sqlite3";
import path from "path";
import fs from "fs";

const costDbPath = process.env.COST_DB_PATH
  ? path.resolve(process.env.COST_DB_PATH)
  : path.join(process.cwd(), "..", "cost.db");

function openCostDb(): Database.Database | null {
  if (!fs.existsSync(costDbPath)) return null;
  const db = new Database(costDbPath, { readonly: true });
  db.pragma("journal_mode = WAL");
  return db;
}

let _db: Database.Database | null | undefined;
function getCostDb(): Database.Database | null {
  if (_db === undefined) {
    _db = openCostDb();
  }
  return _db;
}

export function isCostDbAvailable(): boolean {
  return getCostDb() !== null;
}

export function getTotalCostAllPlayers(): number {
  const db = getCostDb();
  if (!db) return 0;
  const row = db
    .prepare("SELECT COALESCE(SUM(total_cost), 0) as total FROM player_costs")
    .get() as { total: number };
  return row.total;
}

export function getPlayerCosts(): { player_id: string; total_cost: number }[] {
  const db = getCostDb();
  if (!db) return [];
  return db
    .prepare(
      "SELECT player_id, total_cost FROM player_costs ORDER BY total_cost DESC"
    )
    .all() as { player_id: string; total_cost: number }[];
}

export function getCostEventCount(): number {
  const db = getCostDb();
  if (!db) return 0;
  const row = db
    .prepare("SELECT COUNT(*) as count FROM cost_events")
    .get() as { count: number };
  return row.count;
}

export function getDailyCostSplit(
  days: number = 30
): { date: string; premade_cost: number; custom_cost: number }[] {
  const db = getCostDb();
  if (!db) return [];
  return db
    .prepare(
      `SELECT
        date(created_at) as date,
        COALESCE(SUM(CASE WHEN is_custom = 0 THEN amount ELSE 0 END), 0) as premade_cost,
        COALESCE(SUM(CASE WHEN is_custom = 1 THEN amount ELSE 0 END), 0) as custom_cost
      FROM cost_events
      WHERE created_at >= datetime('now', ?)
      GROUP BY date(created_at)
      ORDER BY date`
    )
    .all(`-${days} days`) as {
    date: string;
    premade_cost: number;
    custom_cost: number;
  }[];
}

export function getPerPlayerScenarioBreakdown(
  playerId: string
): {
  scenario_id: string;
  is_custom: number;
  total_cost: number;
  event_count: number;
}[] {
  const db = getCostDb();
  if (!db) return [];
  return db
    .prepare(
      `SELECT scenario_id, is_custom,
        SUM(amount) as total_cost, COUNT(*) as event_count
      FROM cost_events WHERE player_id = ?
      GROUP BY scenario_id, is_custom
      ORDER BY total_cost DESC`
    )
    .all(playerId) as {
    scenario_id: string;
    is_custom: number;
    total_cost: number;
    event_count: number;
  }[];
}

export function getTopScenariosByCost(
  limit: number = 20
): {
  scenario_id: string;
  language: string;
  is_custom: number;
  total_cost: number;
  event_count: number;
}[] {
  const db = getCostDb();
  if (!db) return [];
  return db
    .prepare(
      `SELECT scenario_id, language, is_custom,
        SUM(amount) as total_cost, COUNT(*) as event_count
      FROM cost_events
      GROUP BY scenario_id, language, is_custom
      ORDER BY total_cost DESC
      LIMIT ?`
    )
    .all(limit) as {
    scenario_id: string;
    language: string;
    is_custom: number;
    total_cost: number;
    event_count: number;
  }[];
}

export function getUniquePlayerCount(): number {
  const db = getCostDb();
  if (!db) return 0;
  const row = db
    .prepare("SELECT COUNT(DISTINCT player_id) as count FROM player_costs")
    .get() as { count: number };
  return row.count;
}

export function getCostSplitTotals(): { premade: number; custom: number } {
  const db = getCostDb();
  if (!db) return { premade: 0, custom: 0 };
  const row = db
    .prepare(
      `SELECT
        COALESCE(SUM(CASE WHEN is_custom = 0 THEN amount ELSE 0 END), 0) as premade,
        COALESCE(SUM(CASE WHEN is_custom = 1 THEN amount ELSE 0 END), 0) as custom
      FROM cost_events`
    )
    .get() as { premade: number; custom: number };
  return row;
}

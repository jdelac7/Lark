import Database from "better-sqlite3";
import path from "path";

const dbPath = path.join(process.cwd(), "data", "lark.db");
const db = new Database(dbPath);

// Enable WAL mode for better concurrent read performance
db.pragma("journal_mode = WAL");

// Create tables
db.exec(`
  CREATE TABLE IF NOT EXISTS users (
    id              TEXT PRIMARY KEY,
    email           TEXT UNIQUE NOT NULL,
    password_hash   TEXT,
    name            TEXT,
    polar_customer_id TEXT,
    subscribed      INTEGER DEFAULT 0,
    created_at      TEXT DEFAULT (datetime('now'))
  )
`);

// Migration: allow password_hash to be NULL for OAuth users.
// SQLite can't ALTER COLUMN, but new tables already have it nullable above.
// For existing DBs we just ensure the column exists (it always will).

// Migration: add is_admin column
try {
  db.exec("ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT 0");
} catch {
  // Column already exists — ignore
}

export interface User {
  id: string;
  email: string;
  password_hash: string;
  name: string | null;
  polar_customer_id: string | null;
  subscribed: number;
  is_admin: number;
  created_at: string;
}

export function getUserByEmail(email: string): User | undefined {
  return db.prepare("SELECT * FROM users WHERE email = ?").get(email) as
    | User
    | undefined;
}

export function getUserById(id: string): User | undefined {
  return db.prepare("SELECT * FROM users WHERE id = ?").get(id) as
    | User
    | undefined;
}

export function createUser(
  id: string,
  email: string,
  passwordHash: string,
  name: string | null
): void {
  db.prepare(
    "INSERT INTO users (id, email, password_hash, name) VALUES (?, ?, ?, ?)"
  ).run(id, email, passwordHash, name);
}

export function upsertOAuthUser(
  id: string,
  email: string,
  name: string | null
): User {
  const existing = getUserByEmail(email);
  if (existing) return existing;
  db.prepare(
    "INSERT INTO users (id, email, password_hash, name) VALUES (?, ?, NULL, ?)"
  ).run(id, email, name);
  return getUserById(id)!;
}

export function linkPolarCustomer(
  email: string,
  polarCustomerId: string
): void {
  db.prepare(
    "UPDATE users SET polar_customer_id = ?, subscribed = 1 WHERE email = ?"
  ).run(polarCustomerId, email);
}

export function setSubscriptionStatus(
  polarCustomerId: string,
  active: boolean
): void {
  db.prepare("UPDATE users SET subscribed = ? WHERE polar_customer_id = ?").run(
    active ? 1 : 0,
    polarCustomerId
  );
}

export function isSubscribed(userId: string): boolean {
  const row = db
    .prepare("SELECT subscribed FROM users WHERE id = ?")
    .get(userId) as { subscribed: number } | undefined;
  return row?.subscribed === 1;
}

export function isAdmin(userId: string): boolean {
  const row = db
    .prepare("SELECT is_admin FROM users WHERE id = ?")
    .get(userId) as { is_admin: number } | undefined;
  return row?.is_admin === 1;
}

export function getAllUsers(): User[] {
  return db.prepare("SELECT * FROM users ORDER BY created_at DESC").all() as User[];
}

export default db;

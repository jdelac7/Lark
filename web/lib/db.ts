import Database from "better-sqlite3";
import path from "path";

let _db: Database.Database | null = null;

function getDb(): Database.Database {
  if (_db) return _db;

  const dbPath = path.join(process.cwd(), "data", "lark.db");
  _db = new Database(dbPath);

  // Enable WAL mode for better concurrent read performance
  _db.pragma("journal_mode = WAL");

  // Create tables
  _db.exec(`
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

  // Migration: add is_admin column
  try {
    _db.exec("ALTER TABLE users ADD COLUMN is_admin INTEGER DEFAULT 0");
  } catch {
    // Column already exists — ignore
  }

  // Feedback table
  _db.exec(`
    CREATE TABLE IF NOT EXISTS feedback (
      id         INTEGER PRIMARY KEY AUTOINCREMENT,
      email      TEXT,
      message    TEXT NOT NULL,
      created_at TEXT DEFAULT (datetime('now'))
    )
  `);

  return _db;
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
  return getDb().prepare("SELECT * FROM users WHERE email = ?").get(email) as
    | User
    | undefined;
}

export function getUserById(id: string): User | undefined {
  return getDb().prepare("SELECT * FROM users WHERE id = ?").get(id) as
    | User
    | undefined;
}

export function createUser(
  id: string,
  email: string,
  passwordHash: string,
  name: string | null
): void {
  getDb().prepare(
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
  getDb().prepare(
    "INSERT INTO users (id, email, password_hash, name) VALUES (?, ?, NULL, ?)"
  ).run(id, email, name);
  return getUserById(id)!;
}

export function linkPolarCustomer(
  email: string,
  polarCustomerId: string
): void {
  getDb().prepare(
    "UPDATE users SET polar_customer_id = ?, subscribed = 1 WHERE email = ?"
  ).run(polarCustomerId, email);
}

export function setSubscriptionStatus(
  polarCustomerId: string,
  active: boolean
): void {
  getDb().prepare("UPDATE users SET subscribed = ? WHERE polar_customer_id = ?").run(
    active ? 1 : 0,
    polarCustomerId
  );
}

export function isSubscribed(userId: string): boolean {
  const row = getDb()
    .prepare("SELECT subscribed FROM users WHERE id = ?")
    .get(userId) as { subscribed: number } | undefined;
  return row?.subscribed === 1;
}

export function isAdmin(userId: string): boolean {
  const row = getDb()
    .prepare("SELECT is_admin FROM users WHERE id = ?")
    .get(userId) as { is_admin: number } | undefined;
  return row?.is_admin === 1;
}

export function getAllUsers(): User[] {
  return getDb().prepare("SELECT * FROM users ORDER BY created_at DESC").all() as User[];
}

export interface FeedbackRow {
  id: number;
  email: string | null;
  message: string;
  created_at: string;
}

export function createFeedback(message: string, email: string | null): void {
  getDb()
    .prepare("INSERT INTO feedback (email, message) VALUES (?, ?)")
    .run(email, message);
}

export function getAllFeedback(): FeedbackRow[] {
  return getDb()
    .prepare("SELECT * FROM feedback ORDER BY created_at DESC")
    .all() as FeedbackRow[];
}

export function getFeedbackCount(): number {
  const row = getDb()
    .prepare("SELECT COUNT(*) as count FROM feedback")
    .get() as { count: number };
  return row.count;
}

export default getDb;

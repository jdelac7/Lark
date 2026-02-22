import type { GameMessage, Correction } from "./game-api";

// Mirrors CLI's SaveData structure

export interface SavedSession {
  sessionId: string;
  scenarioName: string;
  lastMessage: GameMessage;
  lastCorrection?: Correction;
  customPrompt?: string;
}

export interface CustomScenario {
  id: string;
  name: string;
  description: string;
}

export interface SaveData {
  completed: Record<string, boolean>;
  sessions: Record<string, SavedSession>;
  customScenarios: CustomScenario[];
}

const SAVE_KEY = "lark-savedata";

function saveKey(scenarioId: string, langCode: string): string {
  return `${scenarioId}:${langCode}`;
}

function newSaveData(): SaveData {
  return { completed: {}, sessions: {}, customScenarios: [] };
}

export function loadSaveData(): SaveData {
  if (typeof window === "undefined") return newSaveData();
  try {
    const raw = localStorage.getItem(SAVE_KEY);
    if (!raw) return newSaveData();
    const sd: SaveData = JSON.parse(raw);
    if (!sd.completed) sd.completed = {};
    if (!sd.sessions) sd.sessions = {};
    if (!sd.customScenarios) sd.customScenarios = [];
    return sd;
  } catch {
    return newSaveData();
  }
}

function persist(sd: SaveData) {
  localStorage.setItem(SAVE_KEY, JSON.stringify(sd));
}

export function isCompleted(scenarioId: string, langCode: string): boolean {
  return loadSaveData().completed[saveKey(scenarioId, langCode)] === true;
}

export function markCompleted(scenarioId: string, langCode: string) {
  const sd = loadSaveData();
  sd.completed[saveKey(scenarioId, langCode)] = true;
  delete sd.sessions[saveKey(scenarioId, langCode)];
  persist(sd);
}

export function saveSession(
  scenarioId: string,
  langCode: string,
  session: SavedSession
) {
  const sd = loadSaveData();
  sd.sessions[saveKey(scenarioId, langCode)] = session;
  persist(sd);
}

export function getSession(
  scenarioId: string,
  langCode: string
): SavedSession | null {
  return loadSaveData().sessions[saveKey(scenarioId, langCode)] || null;
}

export function clearSession(scenarioId: string, langCode: string) {
  const sd = loadSaveData();
  delete sd.sessions[saveKey(scenarioId, langCode)];
  persist(sd);
}

export function getAllSavedSessions(): Record<string, SavedSession> {
  return loadSaveData().sessions;
}

export function addCustomScenario(cs: CustomScenario) {
  const sd = loadSaveData();
  sd.customScenarios.push(cs);
  persist(sd);
}

export function getCustomScenarios(): CustomScenario[] {
  return loadSaveData().customScenarios;
}

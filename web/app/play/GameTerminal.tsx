"use client";

import { useState, useEffect, useRef, useCallback, useMemo } from "react";
import {
  type Scenario,
  type Language,
  type GameMessage,
  type Choice,
  type VocabItem,
  type Correction,
  type Category,
  getScenarios,
  getLanguages,
  startScenarioStream,
  sendChoiceStream,
  sendFreeTextStream,
} from "./game-api";
import {
  isCompleted,
  markCompleted,
  saveSession,
  getSession,
  clearSession,
  addCustomScenario,
  getCustomScenarios,
  type SavedSession,
  type CustomScenario,
} from "./save-state";

type GameState =
  | "selecting_language"
  | "browsing_all_languages"
  | "selecting_category"
  | "selecting_scenario"
  | "entering_custom"
  | "settings"
  | "loading"
  | "streaming"
  | "playing"
  | "finished";

interface Settings {
  showTranslations: boolean;
  showChoices: boolean;
  showVocabulary: boolean;
  showGrammar: boolean;
  explanationLang: string;
}

const defaultSettings: Settings = {
  showTranslations: true,
  showChoices: true,
  showVocabulary: true,
  showGrammar: true,
  explanationLang: "English",
};

const explanationLangOptions = [
  "English", "Español", "Français", "Deutsch",
  "日本語", "中文", "Português", "한국어",
];

type ToggleKey = "showTranslations" | "showChoices" | "showVocabulary" | "showGrammar";

const settingsLabels: { key: ToggleKey; label: string }[] = [
  { key: "showTranslations", label: "Show translations" },
  { key: "showChoices", label: "Show dialog choices" },
  { key: "showVocabulary", label: "Show vocabulary hints" },
  { key: "showGrammar", label: "Show grammar corrections" },
];

function loadSettings(): Settings {
  if (typeof window === "undefined") return defaultSettings;
  try {
    const saved = localStorage.getItem("lark-settings");
    if (saved) return { ...defaultSettings, ...JSON.parse(saved) };
  } catch {}
  return defaultSettings;
}

const POPULAR_CODES = new Set(["es", "fr", "de", "ja", "it", "pt", "ko", "zh"]);

// Extract a JSON string field value from partial/incomplete JSON
function extractField(raw: string, field: string): string {
  const needle = `"${field}":"`;
  const idx = raw.indexOf(needle);
  if (idx < 0) return "";
  const start = idx + needle.length;

  let result = "";
  let escaped = false;
  for (let i = start; i < raw.length; i++) {
    const c = raw[i];
    if (escaped) {
      if (c === "n") result += "\n";
      else if (c === "t") result += "\t";
      else if (c === '"') result += '"';
      else if (c === "\\") result += "\\";
      else if (c === "/") result += "/";
      else result += "\\" + c;
      escaped = false;
      continue;
    }
    if (c === "\\") { escaped = true; continue; }
    if (c === '"') break;
    result += c;
  }
  return result;
}

// Extract completed JSON objects from an array field in partial JSON
function extractArrayItems<T>(raw: string, field: string): T[] {
  const needle = `"${field}":[`;
  const idx = raw.indexOf(needle);
  if (idx < 0) return [];
  const arrStart = idx + needle.length - 1; // points at '['

  // Check if array is complete
  let depth = 0;
  let inStr = false;
  let esc = false;
  let arrEnd = -1;
  for (let i = arrStart; i < raw.length; i++) {
    const c = raw[i];
    if (esc) { esc = false; continue; }
    if (inStr) {
      if (c === "\\") esc = true;
      else if (c === '"') inStr = false;
      continue;
    }
    if (c === '"') inStr = true;
    else if (c === "[") depth++;
    else if (c === "]") {
      depth--;
      if (depth === 0) { arrEnd = i + 1; break; }
    }
  }

  // If array is complete, parse it directly
  if (arrEnd > 0) {
    try { return JSON.parse(raw.slice(arrStart, arrEnd)); } catch {}
  }

  // Array incomplete — extract each complete {...} at depth 1
  const items: T[] = [];
  depth = 0;
  let objStart = -1;
  inStr = false;
  esc = false;
  for (let i = arrStart + 1; i < raw.length; i++) {
    const c = raw[i];
    if (esc) { esc = false; continue; }
    if (inStr) {
      if (c === "\\") esc = true;
      else if (c === '"') inStr = false;
      continue;
    }
    if (c === '"') inStr = true;
    else if (c === "{") {
      if (depth === 0) objStart = i;
      depth++;
    } else if (c === "}") {
      depth--;
      if (depth === 0 && objStart >= 0) {
        try { items.push(JSON.parse(raw.slice(objStart, i + 1))); } catch {}
        objStart = -1;
      }
    }
  }
  return items;
}

export default function GameTerminal() {
  const [state, setState] = useState<GameState>("selecting_language");
  const [prevState, setPrevState] = useState<GameState>("selecting_language");
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [languages, setLanguages] = useState<Language[]>([]);
  const [selectedLanguage, setSelectedLanguage] = useState<Language | null>(null);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [selectedScenario, setSelectedScenario] = useState<Scenario | null>(null);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [message, setMessage] = useState<GameMessage | null>(null);
  const [correction, setCorrection] = useState<Correction | null>(null);
  const [streamTokens, setStreamTokens] = useState("");
  const [freeText, setFreeText] = useState("");
  const [customPrompt, setCustomPrompt] = useState("");
  const [savedCustoms, setSavedCustoms] = useState<CustomScenario[]>([]);
  const [currentScenarioId, setCurrentScenarioId] = useState<string | null>(null);
  const [langSearch, setLangSearch] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [cursor, setCursor] = useState(0);
  const [settings, setSettings] = useState<Settings>(defaultSettings);

  const scrollRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const customInputRef = useRef<HTMLInputElement>(null);
  const langSearchRef = useRef<HTMLInputElement>(null);

  // Load settings and saved custom scenarios on mount
  useEffect(() => {
    setSettings(loadSettings());
    setSavedCustoms(getCustomScenarios());
  }, []);

  // Save settings to localStorage when they change
  function updateSettings(next: Settings) {
    setSettings(next);
    localStorage.setItem("lark-settings", JSON.stringify(next));
  }

  // Parse partial fields from streaming JSON tokens
  const streamParsed = useMemo(() => {
    if (!streamTokens) return null;
    const narrative = extractField(streamTokens, "narrative");
    const translation = extractField(streamTokens, "translation");
    const npcDialog = extractField(streamTokens, "npcDialog");
    const npcDialogTranslation = extractField(streamTokens, "npcDialogTranslation");
    if (!narrative && !npcDialog) return null;
    const choices = extractArrayItems<Choice>(streamTokens, "choices");
    const vocabulary = extractArrayItems<VocabItem>(streamTokens, "vocabulary");
    return { narrative, translation, npcDialog, npcDialogTranslation, choices, vocabulary };
  }, [streamTokens]);

  // Auto-scroll to bottom only during streaming/playing, scroll to top for selection screens
  useEffect(() => {
    if (state === "streaming" || state === "playing" || state === "finished") {
      scrollRef.current?.scrollTo(0, scrollRef.current.scrollHeight);
    } else {
      scrollRef.current?.scrollTo(0, 0);
    }
  }, [state, message, streamTokens, streamParsed]);

  // Scroll active cursor item into view
  useEffect(() => {
    const el = scrollRef.current?.querySelector("[data-active='true']");
    if (el) {
      el.scrollIntoView({ block: "nearest" });
    }
  }, [cursor]);

  // Focus input when playing, entering custom prompt, or searching languages
  useEffect(() => {
    if (state === "playing") inputRef.current?.focus();
    if (state === "entering_custom") customInputRef.current?.focus();
    if (state === "browsing_all_languages") langSearchRef.current?.focus();
  }, [state]);

  const popularLanguages = languages.filter((l) => POPULAR_CODES.has(l.code));
  const filteredLanguages = langSearch
    ? languages.filter((l) =>
        l.name.toLowerCase().includes(langSearch.toLowerCase()) ||
        l.code.toLowerCase().includes(langSearch.toLowerCase())
      )
    : languages;

  // Load scenarios and languages on mount
  useEffect(() => {
    Promise.all([getScenarios(), getLanguages()])
      .then(([s, l]) => {
        setScenarios(s);
        setLanguages(l);
      })
      .catch(() => setError("Failed to connect to game server"));
  }, []);

  // Reset cursor when state changes
  useEffect(() => {
    setCursor(0);
  }, [state]);

  const filteredScenarios = scenarios.filter(
    (s) => s.category === selectedCategory
  );

  // Use ref for keyboard handler so the listener always sees fresh state
  const keyHandlerRef = useRef<(e: KeyboardEvent) => void>(() => {});
  keyHandlerRef.current = (e: KeyboardEvent) => {
    // Escape from screens with focused inputs — handled here with stopImmediatePropagation
    // so the browser can't blur the input or interfere before we navigate back
    if (e.key === "Escape") {
      if (state === "browsing_all_languages") {
        e.preventDefault(); e.stopImmediatePropagation();
        setLangSearch(""); setState("selecting_language"); return;
      }
      if (state === "entering_custom") {
        e.preventDefault(); e.stopImmediatePropagation();
        setState("selecting_category"); return;
      }
    }

    if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

    const num = parseInt(e.key);

    if (state === "selecting_language") {
      const count = popularLanguages.length + 1; // +1 for "All Languages"
      if (!popularLanguages.length) return;
      if (e.key === "ArrowUp") { e.preventDefault(); setCursor((c) => (c - 1 + count) % count); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setCursor((c) => (c + 1) % count); }
      else if (e.key === "Enter") {
        e.preventDefault();
        if (cursor === popularLanguages.length) { setLangSearch(""); setState("browsing_all_languages"); }
        else { selectLanguage(popularLanguages[cursor]); }
      }
      else if (e.key === "s" || e.key === "S") { e.preventDefault(); openSettings(); }
    } else if (state === "browsing_all_languages") {
      // Handle keys when search input lost focus (clicked elsewhere)
      if (e.key === "Escape") { e.preventDefault(); setLangSearch(""); setState("selecting_language"); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setCursor((c) => Math.min(c + 1, filteredLanguages.length - 1)); }
      else if (e.key === "ArrowUp") { e.preventDefault(); setCursor((c) => Math.max(c - 1, 0)); }
      else if (e.key === "Enter") { e.preventDefault(); if (filteredLanguages.length > 0) selectLanguage(filteredLanguages[cursor]); }
      else { langSearchRef.current?.focus(); }
    } else if (state === "selecting_category") {
      const count = 3;
      if (e.key === "ArrowUp") { e.preventDefault(); setCursor((c) => (c - 1 + count) % count); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setCursor((c) => (c + 1) % count); }
      else if (e.key === "Enter") {
        e.preventDefault();
        if (cursor === 2) { setState("entering_custom"); }
        else { selectCategory(cursor === 0 ? "everyday" : "adventure"); }
      }
      else if (e.key === "Escape") { e.preventDefault(); setState("selecting_language"); }
      else if (num === 1) { e.preventDefault(); selectCategory("everyday"); }
      else if (num === 2) { e.preventDefault(); selectCategory("adventure"); }
      else if (num === 3) { e.preventDefault(); setState("entering_custom"); }
    } else if (state === "selecting_scenario") {
      const count = filteredScenarios.length;
      if (!count) return;
      if (e.key === "ArrowUp") { e.preventDefault(); setCursor((c) => (c - 1 + count) % count); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setCursor((c) => (c + 1) % count); }
      else if (e.key === "Enter") { e.preventDefault(); selectScenario(filteredScenarios[cursor]); }
      else if (e.key === "Escape") { e.preventDefault(); setState("selecting_category"); }
      else if (num >= 1 && num <= count) { e.preventDefault(); selectScenario(filteredScenarios[num - 1]); }
    } else if (state === "settings") {
      const count = settingsLabels.length + 1; // +1 for explanation language row
      if (e.key === "ArrowUp") { e.preventDefault(); setCursor((c) => (c - 1 + count) % count); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setCursor((c) => (c + 1) % count); }
      else if (e.key === "Enter") {
        e.preventDefault();
        if (cursor < settingsLabels.length) {
          const key = settingsLabels[cursor].key;
          updateSettings({ ...settings, [key]: !settings[key] });
        } else {
          // Cycle explanation language
          const curIdx = explanationLangOptions.indexOf(settings.explanationLang);
          const nextLang = explanationLangOptions[(curIdx + 1) % explanationLangOptions.length];
          updateSettings({ ...settings, explanationLang: nextLang });
        }
      }
      else if (e.key === "Escape") { e.preventDefault(); setState(prevState); }
    } else if (state === "playing" && message?.choices?.length) {
      if (num >= 1 && num <= message.choices.length) {
        handleChoice(num - 1);
      }
    }
  };

  // Attach a single stable listener in CAPTURE phase so it fires before
  // the browser's default input handling (which blurs on Escape)
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      keyHandlerRef.current(e);
    }
    window.addEventListener("keydown", handleKeyDown, true);
    return () => window.removeEventListener("keydown", handleKeyDown, true);
  }, []);

  const handleToken = useCallback((token: string) => {
    setStreamTokens((prev) => prev + token);
  }, []);

  function openSettings() {
    setPrevState(state);
    setState("settings");
  }

  function selectLanguage(lang: Language) {
    setSelectedLanguage(lang);
    setState("selecting_category");
  }

  function selectCategory(cat: Category) {
    setSelectedCategory(cat);
    setState("selecting_scenario");
  }

  async function selectScenario(scenario: Scenario) {
    setSelectedScenario(scenario);
    setCurrentScenarioId(scenario.id);
    const langCode = selectedLanguage!.code;

    // Check for saved session to resume
    const saved = getSession(scenario.id, langCode);
    if (saved) {
      setSessionId(saved.sessionId);
      setMessage(saved.lastMessage);
      setCorrection(saved.lastCorrection || null);
      setState(saved.lastMessage.finished ? "finished" : "playing");
      return;
    }

    setState("streaming");
    setStreamTokens("");
    setError(null);

    try {
      const result = await startScenarioStream(
        scenario.id,
        langCode,
        handleToken,
        undefined,
        settings.explanationLang
      );
      setSessionId(result.sessionId);
      setMessage(result.message);
      setStreamTokens("");
      const newState = result.message.finished ? "finished" : "playing";
      setState(newState);

      // Save session for resume
      saveSession(scenario.id, langCode, {
        sessionId: result.sessionId,
        scenarioName: scenario.name,
        lastMessage: result.message,
      });
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to start scenario"
      );
      setState("selecting_scenario");
    }
  }

  async function startCustom(e: React.FormEvent) {
    e.preventDefault();
    if (!customPrompt.trim()) return;
    const langCode = selectedLanguage!.code;
    const prompt = customPrompt.trim();

    // Generate a unique ID and save the custom scenario
    const customId = `custom_${Date.now().toString(36)}`;
    const cs: CustomScenario = {
      id: customId,
      name: prompt.length > 40 ? prompt.slice(0, 40) + "..." : prompt,
      description: prompt,
    };
    addCustomScenario(cs);
    setSavedCustoms(getCustomScenarios());
    setCurrentScenarioId(customId);

    setState("streaming");
    setStreamTokens("");
    setError(null);

    try {
      const result = await startScenarioStream(
        "custom",
        langCode,
        handleToken,
        prompt,
        settings.explanationLang
      );
      setSessionId(result.sessionId);
      setMessage(result.message);
      setStreamTokens("");
      const newState = result.message.finished ? "finished" : "playing";
      setState(newState);

      saveSession(customId, langCode, {
        sessionId: result.sessionId,
        scenarioName: cs.name,
        lastMessage: result.message,
        customPrompt: prompt,
      });
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to start custom scenario"
      );
      setState("entering_custom");
    }
  }

  async function resumeCustom(cs: CustomScenario) {
    const langCode = selectedLanguage!.code;
    const saved = getSession(cs.id, langCode);
    if (saved) {
      setCurrentScenarioId(cs.id);
      setSessionId(saved.sessionId);
      setMessage(saved.lastMessage);
      setCorrection(saved.lastCorrection || null);
      setState(saved.lastMessage.finished ? "finished" : "playing");
      return;
    }
    // No saved session — start fresh with the custom prompt
    setCurrentScenarioId(cs.id);
    setState("streaming");
    setStreamTokens("");
    setError(null);

    try {
      const result = await startScenarioStream(
        "custom",
        langCode,
        handleToken,
        cs.description,
        settings.explanationLang
      );
      setSessionId(result.sessionId);
      setMessage(result.message);
      setStreamTokens("");
      setState(result.message.finished ? "finished" : "playing");
      saveSession(cs.id, langCode, {
        sessionId: result.sessionId,
        scenarioName: cs.name,
        lastMessage: result.message,
        customPrompt: cs.description,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to start scenario");
      setState("entering_custom");
    }
  }

  function afterTurn(msg: GameMessage, corr: Correction | null) {
    const scenId = currentScenarioId;
    const langCode = selectedLanguage?.code;
    if (!scenId || !langCode || !sessionId) return;

    if (msg.finished) {
      markCompleted(scenId, langCode);
      clearSession(scenId, langCode);
    } else {
      saveSession(scenId, langCode, {
        sessionId: sessionId,
        scenarioName: selectedScenario?.name || "Custom",
        lastMessage: msg,
        lastCorrection: corr || undefined,
      });
    }
  }

  async function handleChoice(index: number) {
    if (!sessionId || state !== "playing") return;
    setState("streaming");
    setStreamTokens("");
    setCorrection(null);

    try {
      const result = await sendChoiceStream(sessionId, index, handleToken);
      setMessage(result.message);
      setCorrection(result.correction || null);
      setStreamTokens("");
      setState(result.message.finished ? "finished" : "playing");
      afterTurn(result.message, result.correction || null);
    } catch (err) {
      handleSessionError(err, "Failed to send choice");
    }
  }

  async function handleFreeText(e: React.FormEvent) {
    e.preventDefault();
    if (!sessionId || !freeText.trim() || state !== "playing") return;

    const text = freeText.trim();
    setFreeText("");

    // If user typed a number matching a choice, treat it as a choice selection
    const num = parseInt(text);
    if (text === String(num) && message?.choices && num >= 1 && num <= message.choices.length) {
      handleChoice(num - 1);
      return;
    }
    setState("streaming");
    setStreamTokens("");
    setCorrection(null);

    try {
      const result = await sendFreeTextStream(sessionId, text, handleToken);
      setMessage(result.message);
      setCorrection(result.correction || null);
      setStreamTokens("");
      setState(result.message.finished ? "finished" : "playing");
      afterTurn(result.message, result.correction || null);
    } catch (err) {
      handleSessionError(err, "Failed to send input");
    }
  }

  async function handleSessionError(err: unknown, fallbackMsg: string) {
    // Session lost on server — clear stale save and restart the scenario automatically
    if (currentScenarioId && selectedLanguage) {
      clearSession(currentScenarioId, selectedLanguage.code);
    }
    setSessionId(null);

    // Try to restart the scenario transparently
    if (selectedScenario && selectedLanguage) {
      setError(null);
      setState("streaming");
      setStreamTokens("");
      try {
        const result = await startScenarioStream(
          selectedScenario.id,
          selectedLanguage.code,
          handleToken,
          undefined,
          settings.explanationLang
        );
        setSessionId(result.sessionId);
        setMessage(result.message);
        setStreamTokens("");
        setState(result.message.finished ? "finished" : "playing");
        if (currentScenarioId) {
          saveSession(currentScenarioId, selectedLanguage.code, {
            sessionId: result.sessionId,
            scenarioName: selectedScenario.name,
            lastMessage: result.message,
          });
        }
        return;
      } catch {
        // Restart also failed — fall back to scenario selection
      }
    }

    setError("Connection lost — please select a scenario");
    setState("selecting_scenario");
  }

  function resetGame() {
    setState("selecting_language");
    setSelectedLanguage(null);
    setSelectedCategory(null);
    setSelectedScenario(null);
    setCurrentScenarioId(null);
    setSessionId(null);
    setMessage(null);
    setCorrection(null);
    setStreamTokens("");
    setFreeText("");
    setCustomPrompt("");
    setLangSearch("");
    setError(null);
  }

  const cursorPrefix = (active: boolean) =>
    active ? "text-accent" : "text-text-dim";

  return (
    <div className="mx-auto flex h-[calc(100vh-4rem)] max-w-4xl flex-col px-4 pt-16">
      {/* Terminal window */}
      <div className="flex flex-1 flex-col overflow-hidden border border-accent/30 bg-bg-card">
        {/* Title bar */}
        <div className="flex items-center gap-2 border-b border-border px-4 py-2">
          <a href="/" className="h-2.5 w-2.5 rounded-full bg-accent/60 transition-colors hover:bg-accent" title="Back to home" />
          <div className="h-2.5 w-2.5 rounded-full bg-yellow/60" />
          <div className="h-2.5 w-2.5 rounded-full bg-text-dim/40" />
          <span className="ml-2 text-xs text-text-dim">
            lark
            {state === "settings"
              ? " — settings"
              : selectedScenario
                ? ` — ${selectedScenario.name} · ${selectedLanguage?.name}`
                : " — play"}
          </span>
        </div>

        {/* Content area */}
        <div
          ref={scrollRef}
          className="flex-1 overflow-y-auto p-4 font-mono text-sm"
        >
          {error && (
            <div className="mb-4 border border-yellow/30 bg-yellow/10 px-3 py-2 text-yellow">
              ERROR: {error}
            </div>
          )}

          {/* Language selection — popular */}
          {state === "selecting_language" && (
            <div>
              <div className="mb-4 text-text-dim">$ lark play</div>
              <div className="mb-4 text-accent">
                SELECT YOUR TARGET LANGUAGE:
              </div>
              {popularLanguages.length === 0 && !error ? (
                <div className="space-y-2">
                  <span className="text-text-dim">Connecting to game server</span>
                  <span className="animate-pulse text-accent">...</span>
                </div>
              ) : (
              <div className="space-y-1">
                {popularLanguages.map((lang, i) => (
                  <button
                    key={lang.code}
                    data-active={cursor === i} tabIndex={-1}
                    onClick={() => selectLanguage(lang)}
                    className={`block w-full text-left transition-colors hover:text-accent ${cursor === i ? "text-accent" : ""}`}
                  >
                    <span className={cursorPrefix(cursor === i)}>
                      {cursor === i ? ">" : " "} {i + 1}.
                    </span>{" "}
                    <span className={cursor === i ? "text-accent" : "text-cyan"}>
                      {lang.name}
                    </span>{" "}
                    <span className="text-text-dim">({lang.code})</span>
                  </button>
                ))}
                <button
                  data-active={cursor === popularLanguages.length} tabIndex={-1}
                  onClick={() => { setLangSearch(""); setState("browsing_all_languages"); }}
                  className={`block w-full text-left transition-colors hover:text-accent ${cursor === popularLanguages.length ? "text-accent" : ""}`}
                >
                  <span className={cursorPrefix(cursor === popularLanguages.length)}>
                    {cursor === popularLanguages.length ? ">" : " "} {popularLanguages.length + 1}.
                  </span>{" "}
                  <span className={cursor === popularLanguages.length ? "text-accent" : "text-purple"}>
                    All Languages
                  </span>{" "}
                  <span className="text-text-dim">({languages.length} available)</span>
                </button>
              </div>
              )}
            </div>
          )}

          {/* All languages — scrollable with search */}
          {state === "browsing_all_languages" && (
            <div>
              <div className="mb-4 text-accent">ALL LANGUAGES</div>
              <div className="mb-3 flex items-center">
                <span className="mr-2 text-accent">search:</span>
                <input
                  ref={langSearchRef}
                  type="text"
                  value={langSearch}
                  onChange={(e) => { setLangSearch(e.target.value); setCursor(0); }}
                  onKeyDown={(e) => {
                    if (e.key === "ArrowDown") {
                      e.preventDefault();
                      setCursor((c) => Math.min(c + 1, filteredLanguages.length - 1));
                    } else if (e.key === "ArrowUp") {
                      e.preventDefault();
                      setCursor((c) => Math.max(c - 1, 0));
                    } else if (e.key === "Enter") {
                      e.preventDefault();
                      if (filteredLanguages.length > 0) {
                        selectLanguage(filteredLanguages[cursor]);
                      }
                    }
                  }}
                  placeholder="Type to filter..."
                  className="flex-1 bg-transparent text-text outline-none placeholder:text-text-dim/50"
                />
              </div>
              <div className="border-t border-border pt-3">
                <div className="space-y-1">
                  {filteredLanguages.map((lang, i) => (
                    <button
                      key={lang.code}
                      tabIndex={-1}
                      data-active={cursor === i}
                      onClick={() => selectLanguage(lang)}
                      className={`block w-full text-left transition-colors hover:text-accent ${cursor === i ? "text-accent" : ""}`}
                    >
                      <span className={cursorPrefix(cursor === i)}>
                        {cursor === i ? ">" : " "}
                      </span>{" "}
                      <span className={cursor === i ? "text-accent" : "text-cyan"}>{lang.name}</span>{" "}
                      <span className="text-text-dim">({lang.code})</span>
                    </button>
                  ))}
                  {filteredLanguages.length === 0 && (
                    <div className="text-text-dim">No languages match your search.</div>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Category selection */}
          {state === "selecting_category" && (
            <div>
              <div className="mb-4 text-text-dim">
                Language: {selectedLanguage?.name}
              </div>
              <div className="mb-4 text-accent">SELECT CATEGORY:</div>
              <div className="space-y-2">
                {([
                  { key: "everyday", label: "Everyday", desc: "Real-world situations" },
                  { key: "adventure", label: "Adventure", desc: "Fantasy & sci-fi scenarios" },
                  { key: "custom", label: "Custom", desc: "Describe your own scenario" },
                ]).map((item, i) => (
                  <button
                    key={item.key}
                    data-active={cursor === i} tabIndex={-1}
                    onClick={() => {
                      if (item.key === "custom") setState("entering_custom");
                      else selectCategory(item.key as Category);
                    }}
                    className={`block w-full text-left transition-colors hover:text-accent ${cursor === i ? "text-accent" : ""}`}
                  >
                    <span className={cursorPrefix(cursor === i)}>
                      {cursor === i ? ">" : " "} {i + 1}.
                    </span>{" "}
                    <span className={cursor === i ? "text-accent" : "text-cyan"}>
                      {item.label}
                    </span>
                    <span className="ml-2 text-text-dim">— {item.desc}</span>
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Scenario selection */}
          {state === "selecting_scenario" && (
            <div>
              <div className="mb-2 text-text-dim">
                {selectedLanguage?.name} · {selectedCategory}
              </div>
              <div className="mb-4 text-accent">SELECT SCENARIO:</div>
              <div className="space-y-2">
                {filteredScenarios.map((s, i) => {
                  const done = selectedLanguage && isCompleted(s.id, selectedLanguage.code);
                  const hasSave = selectedLanguage && getSession(s.id, selectedLanguage.code);
                  return (
                    <button
                      key={s.id}
                      data-active={cursor === i} tabIndex={-1}
                      onClick={() => selectScenario(s)}
                      className={`block w-full text-left transition-colors hover:text-accent ${cursor === i ? "text-accent" : ""}`}
                    >
                      <span className={cursorPrefix(cursor === i)}>
                        {cursor === i ? ">" : " "} {i + 1}.
                      </span>{" "}
                      <span className={cursor === i ? "text-accent" : "text-cyan"}>
                        {s.name}
                      </span>
                      <span className="ml-2 text-text-dim">[{s.difficulty}]</span>
                      {done && <span className="ml-2 text-green">&#10003;</span>}
                      {!done && hasSave && <span className="ml-2 text-yellow">[resume]</span>}
                      <div className="ml-6 text-xs text-text-dim">
                        {s.description}
                      </div>
                    </button>
                  );
                })}
              </div>
            </div>
          )}

          {/* Custom prompt input */}
          {state === "entering_custom" && (
            <div>
              <div className="mb-4 text-text-dim">
                Language: {selectedLanguage?.name}
              </div>
              <div className="mb-4 text-accent">CUSTOM SCENARIO</div>

              {/* Previously saved custom scenarios */}
              {savedCustoms.length > 0 && (
                <div className="mb-4">
                  <div className="mb-2 text-xs text-text-dim">PREVIOUS CUSTOM SCENARIOS:</div>
                  <div className="space-y-1">
                    {savedCustoms.map((cs) => {
                      const done = selectedLanguage && isCompleted(cs.id, selectedLanguage.code);
                      const hasSave = selectedLanguage && getSession(cs.id, selectedLanguage.code);
                      return (
                        <button
                          key={cs.id}
                          tabIndex={-1}
                          onClick={() => resumeCustom(cs)}
                          className="block w-full text-left text-sm transition-colors hover:text-accent"
                        >
                          <span className="text-text-dim">&middot;</span>{" "}
                          <span className="text-cyan">{cs.name}</span>
                          {done && <span className="ml-2 text-green">&#10003;</span>}
                          {!done && hasSave && <span className="ml-2 text-yellow">[resume]</span>}
                        </button>
                      );
                    })}
                  </div>
                  <div className="my-3 border-t border-border" />
                </div>
              )}

              <div className="mb-4 text-text-dim">
                Describe a new scenario:
              </div>
              <form onSubmit={startCustom} className="flex items-center">
                <span className="mr-2 text-accent">$</span>
                <input
                  ref={customInputRef}
                  type="text"
                  value={customPrompt}
                  onChange={(e) => setCustomPrompt(e.target.value)}
                  placeholder="e.g. Ordering street food in a night market..."
                  className="flex-1 bg-transparent text-text outline-none placeholder:text-text-dim/50"
                />
              </form>
            </div>
          )}

          {/* Settings */}
          {state === "settings" && (
            <div>
              <div className="mb-4 text-text-dim">$ lark settings</div>
              <div className="mb-4 text-accent">SETTINGS</div>
              <div className="space-y-2">
                {settingsLabels.map((s, i) => (
                  <button
                    key={s.key}
                    data-active={cursor === i} tabIndex={-1}
                    onClick={() =>
                      updateSettings({ ...settings, [s.key]: !settings[s.key] })
                    }
                    className={`block w-full text-left transition-colors hover:text-accent ${cursor === i ? "text-accent" : ""}`}
                  >
                    <span className={cursorPrefix(cursor === i)}>
                      {cursor === i ? ">" : " "}
                    </span>{" "}
                    <span className={cursor === i ? "text-accent" : "text-text"}>
                      {s.label}
                    </span>
                    <span className={`ml-2 font-bold ${settings[s.key] ? "text-green" : "text-text-dim"}`}>
                      [{settings[s.key] ? "ON" : "OFF"}]
                    </span>
                  </button>
                ))}
                {/* Explanation language cycle picker */}
                <button
                  data-active={cursor === settingsLabels.length} tabIndex={-1}
                  onClick={() => {
                    const curIdx = explanationLangOptions.indexOf(settings.explanationLang);
                    const nextLang = explanationLangOptions[(curIdx + 1) % explanationLangOptions.length];
                    updateSettings({ ...settings, explanationLang: nextLang });
                  }}
                  className={`block w-full text-left transition-colors hover:text-accent ${cursor === settingsLabels.length ? "text-accent" : ""}`}
                >
                  <span className={cursorPrefix(cursor === settingsLabels.length)}>
                    {cursor === settingsLabels.length ? ">" : " "}
                  </span>{" "}
                  <span className={cursor === settingsLabels.length ? "text-accent" : "text-text"}>
                    Explanation language
                  </span>
                  <span className="ml-2 font-bold text-cyan">
                    [{settings.explanationLang}]
                  </span>
                </button>
              </div>
            </div>
          )}

          {/* Loading */}
          {state === "loading" && (
            <div className="text-accent">
              Connecting...
              <span className="terminal-cursor" />
            </div>
          )}

          {/* Streaming — show parsed fields as they arrive */}
          {state === "streaming" && (
            <div className="space-y-4">
              {streamParsed ? (
                <>
                  {streamParsed.narrative && (
                    <div>
                      <div className="whitespace-pre-wrap text-green">
                        {streamParsed.narrative}
                        <span className="terminal-cursor" />
                      </div>
                      {settings.showTranslations && streamParsed.translation && (
                        <div className="mt-1 whitespace-pre-wrap text-text-dim italic">
                          {streamParsed.translation}
                        </div>
                      )}
                    </div>
                  )}
                  {streamParsed.npcDialog && (
                    <div className="border-l-2 border-cyan/40 pl-3">
                      <div className="text-cyan">
                        &ldquo;{streamParsed.npcDialog}&rdquo;
                      </div>
                      {settings.showTranslations && streamParsed.npcDialogTranslation && (
                        <div className="mt-1 text-text-dim italic">
                          &ldquo;{streamParsed.npcDialogTranslation}&rdquo;
                        </div>
                      )}
                    </div>
                  )}
                  {settings.showVocabulary && streamParsed.vocabulary.length > 0 && (
                    <div className="border-t border-border pt-3">
                      <div className="mb-2 text-xs font-bold text-purple">
                        VOCABULARY
                      </div>
                      <div className="space-y-1">
                        {streamParsed.vocabulary.map((v, i) => (
                          <div key={i} className="text-xs">
                            <span className="text-purple">{v.word}</span>
                            <span className="text-text-dim"> — </span>
                            <span className="text-text">{v.translation}</span>
                            {v.usage && (
                              <span className="ml-2 text-text-dim italic">
                                ({v.usage})
                              </span>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                  {settings.showChoices && streamParsed.choices.length > 0 && (
                    <div className="border-t border-border pt-3">
                      <div className="space-y-1">
                        {streamParsed.choices.map((c, i) => (
                          <div key={i} className="text-left">
                            <span className="text-text-dim">{i + 1}.</span>{" "}
                            <span className="text-cyan">{c.text}</span>
                            {settings.showTranslations && c.translation && (
                              <span className="ml-2 text-text-dim italic">
                                ({c.translation})
                              </span>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </>
              ) : (
                <div className="text-accent">
                  <span className="terminal-cursor" />
                </div>
              )}
            </div>
          )}

          {/* Playing - main game screen */}
          {(state === "playing" || state === "finished") && message && (
            <div className="space-y-4">
              {/* Narrative */}
              <div>
                <div className="whitespace-pre-wrap text-green">
                  {message.narrative}
                </div>
                {settings.showTranslations && message.translation && (
                  <div className="mt-1 whitespace-pre-wrap text-text-dim italic">
                    {message.translation}
                  </div>
                )}
              </div>

              {/* NPC Dialog */}
              {message.npcDialog && (
                <div className="border-l-2 border-cyan/40 pl-3">
                  <div className="text-cyan">
                    &ldquo;{message.npcDialog}&rdquo;
                  </div>
                  {settings.showTranslations && message.npcDialogTranslation && (
                    <div className="mt-1 text-text-dim italic">
                      &ldquo;{message.npcDialogTranslation}&rdquo;
                    </div>
                  )}
                </div>
              )}

              {/* Correction */}
              {settings.showGrammar && correction && (
                <div className="border border-yellow/30 bg-yellow/5 p-3">
                  <div className="mb-1 text-xs font-bold text-yellow">
                    CORRECTION
                  </div>
                  <div className="text-text-dim line-through">
                    {correction.original}
                  </div>
                  <div className="text-green">{correction.corrected}</div>
                  <div className="mt-1 text-xs text-yellow">
                    {correction.explanation}
                  </div>
                </div>
              )}

              {/* Vocabulary */}
              {settings.showVocabulary &&
                message.vocabulary &&
                message.vocabulary.length > 0 && (
                  <div className="border-t border-border pt-3">
                    <div className="mb-2 text-xs font-bold text-purple">
                      VOCABULARY
                    </div>
                    <div className="space-y-1">
                      {message.vocabulary.map((v, i) => (
                        <div key={i} className="text-xs">
                          <span className="text-purple">{v.word}</span>
                          <span className="text-text-dim"> — </span>
                          <span className="text-text">{v.translation}</span>
                          {v.usage && (
                            <span className="ml-2 text-text-dim italic">
                              ({v.usage})
                            </span>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

              {/* Finished */}
              {state === "finished" && (
                <div className="border-t border-accent/30 pt-4">
                  <div className="mb-3 text-accent font-bold">
                    SCENARIO COMPLETE
                  </div>
                  <button
                    onClick={resetGame}
                    className="border border-accent px-4 py-2 text-sm font-bold text-accent transition-colors hover:bg-accent/10"
                  >
                    &gt; PLAY AGAIN
                  </button>
                </div>
              )}

              {/* Choices */}
              {state === "playing" &&
                settings.showChoices &&
                message.choices &&
                message.choices.length > 0 && (
                  <div className="border-t border-border pt-3">
                    <div className="space-y-1">
                      {message.choices.map((c, i) => (
                        <button
                          key={i}
                          tabIndex={-1}
                          onClick={() => handleChoice(i)}
                          className="block w-full text-left transition-colors hover:text-accent"
                        >
                          <span className="text-text-dim">{i + 1}.</span>{" "}
                          <span className="text-cyan">{c.text}</span>
                          {settings.showTranslations && (
                            <span className="ml-2 text-text-dim italic">
                              ({c.translation})
                            </span>
                          )}
                        </button>
                      ))}
                    </div>
                  </div>
                )}
            </div>
          )}
        </div>

        {/* Navigation hints footer */}
        {(state === "selecting_language" ||
          state === "browsing_all_languages" ||
          state === "selecting_category" ||
          state === "selecting_scenario" ||
          state === "entering_custom" ||
          state === "settings") && (
          <div className="border-t border-border px-4 py-2 text-xs text-text-dim font-mono">
            {state === "selecting_language" && "[↑↓] navigate · [Enter] select · [S] settings"}
            {state === "browsing_all_languages" && "[↑↓] navigate · [Enter] select · [Esc] back"}
            {state === "selecting_category" && "[↑↓] navigate · [Enter] select · [Esc] back"}
            {state === "selecting_scenario" && "[↑↓] navigate · [Enter] select · [Esc] back"}
            {state === "entering_custom" && "[Enter] start · [Esc] back"}
            {state === "settings" && "[↑↓] navigate · [Enter] toggle · [Esc] back"}
          </div>
        )}

        {/* Input bar */}
        {state === "playing" && !message?.finished && (
          <form
            onSubmit={handleFreeText}
            className="flex items-center border-t border-border px-4 py-3"
          >
            <span className="mr-2 text-accent">$</span>
            <input
              ref={inputRef}
              type="text"
              value={freeText}
              onChange={(e) => setFreeText(e.target.value)}
              placeholder="Type your own response..."
              className="flex-1 bg-transparent font-mono text-sm text-text outline-none placeholder:text-text-dim/50"
            />
            <button
              type="submit"
              className="ml-2 text-xs text-accent transition-colors hover:text-accent-dim"
            >
              [send]
            </button>
          </form>
        )}
      </div>
    </div>
  );
}

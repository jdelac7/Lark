const SERVER_URL =
  process.env.NEXT_PUBLIC_GAME_SERVER_URL || "http://localhost:9292";

// Types mirroring Go api/types.go

export type InputMode = "choice" | "free_text";
export type Difficulty = "beginner" | "intermediate" | "advanced";
export type Category = "everyday" | "adventure";

export interface Scenario {
  id: string;
  name: string;
  description: string;
  difficulty: Difficulty;
  category: Category;
}

export interface Language {
  code: string;
  name: string;
}

export interface GameMessage {
  narrative: string;
  translation: string;
  npcDialog?: string;
  npcDialogTranslation?: string;
  choices: Choice[];
  vocabulary: VocabItem[];
  finished: boolean;
}

export interface Choice {
  text: string;
  translation: string;
}

export interface VocabItem {
  word: string;
  translation: string;
  usage?: string;
}

export interface Correction {
  original: string;
  corrected: string;
  explanation: string;
}

export interface StartResponse {
  sessionId: string;
  message: GameMessage;
}

export interface PlayerInputResponse {
  message: GameMessage;
  correction?: Correction;
}

// SSE stream event types
interface StreamToken {
  token: string;
}

interface StreamDone {
  done: true;
  sessionId?: string;
  message: GameMessage;
  correction?: Correction;
}

interface StreamError {
  error: string;
}

type StreamEvent = StreamToken | StreamDone | StreamError;

// Fetch helpers

export async function getScenarios(): Promise<Scenario[]> {
  const res = await fetch(`${SERVER_URL}/api/v1/scenarios`);
  if (!res.ok) throw new Error("Failed to fetch scenarios");
  return res.json();
}

export async function getLanguages(): Promise<Language[]> {
  const res = await fetch(`${SERVER_URL}/api/v1/languages`);
  if (!res.ok) throw new Error("Failed to fetch languages");
  return res.json();
}

export async function startScenarioStream(
  scenarioId: string,
  language: string,
  onToken: (token: string) => void,
  customPrompt?: string,
  explanationLang?: string
): Promise<{ sessionId: string; message: GameMessage }> {
  const res = await fetch(`${SERVER_URL}/api/v1/scenarios/start/stream`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      scenarioId,
      language,
      ...(customPrompt ? { customPrompt } : {}),
      ...(explanationLang ? { explanationLang } : {}),
    }),
  });

  if (!res.ok) throw new Error("Failed to start scenario");

  return parseSSEStream<{ sessionId: string; message: GameMessage }>(
    res,
    onToken,
    (event: StreamDone) => ({
      sessionId: event.sessionId!,
      message: event.message,
    })
  );
}

export async function sendChoiceStream(
  sessionId: string,
  choiceIndex: number,
  onToken: (token: string) => void
): Promise<PlayerInputResponse> {
  const res = await fetch(`${SERVER_URL}/api/v1/game/input/stream`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      sessionId,
      mode: "choice" as InputMode,
      choiceIndex,
    }),
  });

  if (!res.ok) throw new Error("Failed to send choice");

  return parseSSEStream<PlayerInputResponse>(
    res,
    onToken,
    (event: StreamDone) => ({
      message: event.message,
      correction: event.correction,
    })
  );
}

export async function sendFreeTextStream(
  sessionId: string,
  text: string,
  onToken: (token: string) => void
): Promise<PlayerInputResponse> {
  const res = await fetch(`${SERVER_URL}/api/v1/game/input/stream`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      sessionId,
      mode: "free_text" as InputMode,
      text,
    }),
  });

  if (!res.ok) throw new Error("Failed to send input");

  return parseSSEStream<PlayerInputResponse>(
    res,
    onToken,
    (event: StreamDone) => ({
      message: event.message,
      correction: event.correction,
    })
  );
}

async function parseSSEStream<T>(
  response: Response,
  onToken: (token: string) => void,
  extractResult: (event: StreamDone) => T
): Promise<T> {
  const reader = response.body!.getReader();
  const decoder = new TextDecoder();
  let buffer = "";

  return new Promise<T>(async (resolve, reject) => {
    try {
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split("\n");
        buffer = lines.pop() || "";

        for (const line of lines) {
          if (!line.startsWith("data: ")) continue;
          const jsonStr = line.slice(6).trim();
          if (!jsonStr) continue;

          try {
            const event: StreamEvent = JSON.parse(jsonStr);

            if ("error" in event) {
              reject(new Error(event.error));
              return;
            }

            if ("done" in event && event.done) {
              resolve(extractResult(event));
              return;
            }

            if ("token" in event) {
              onToken(event.token);
            }
          } catch {
            // Skip malformed JSON lines
          }
        }
      }

      reject(new Error("Stream ended without done event"));
    } catch (err) {
      reject(err);
    }
  });
}

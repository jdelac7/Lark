const GAME_SERVER =
  process.env.GAME_SERVER_INTERNAL_URL || "http://localhost:9292";

export async function POST(request: Request) {
  const body = await request.text();

  const upstream = await fetch(`${GAME_SERVER}/api/v1/game/input/stream`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body,
  });

  if (!upstream.ok) {
    return new Response(await upstream.text(), { status: upstream.status });
  }

  // Pipe the SSE stream through without buffering
  return new Response(upstream.body, {
    headers: {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
    },
  });
}

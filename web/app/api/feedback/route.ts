import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { createFeedback } from "@/lib/db";

export async function POST(request: NextRequest) {
  const { message } = await request.json();

  if (!message || typeof message !== "string" || message.trim().length === 0) {
    return NextResponse.json(
      { error: "Message is required" },
      { status: 400 }
    );
  }

  if (message.length > 2000) {
    return NextResponse.json(
      { error: "Message must be under 2000 characters" },
      { status: 400 }
    );
  }

  // Capture email if logged in, otherwise null
  let email: string | null = null;
  try {
    const session = await auth();
    email = session?.user?.email ?? null;
  } catch {
    // Not logged in — that's fine
  }

  createFeedback(message.trim(), email);

  return NextResponse.json({ success: true });
}

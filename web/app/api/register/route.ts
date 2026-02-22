import { NextRequest, NextResponse } from "next/server";
import { hash } from "bcryptjs";
import { nanoid } from "nanoid";
import { createUser, getUserByEmail } from "@/lib/db";

export async function POST(request: NextRequest) {
  const { email, password, name } = await request.json();

  if (!email || !password) {
    return NextResponse.json(
      { error: "Email and password are required" },
      { status: 400 }
    );
  }

  if (password.length < 6) {
    return NextResponse.json(
      { error: "Password must be at least 6 characters" },
      { status: 400 }
    );
  }

  const existing = getUserByEmail(email);
  if (existing) {
    return NextResponse.json(
      { error: "An account with this email already exists" },
      { status: 409 }
    );
  }

  const id = nanoid();
  const passwordHash = await hash(password, 12);
  createUser(id, email, passwordHash, name || null);

  return NextResponse.json({ success: true });
}

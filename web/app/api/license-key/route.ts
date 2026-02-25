import { NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { polar } from "@/lib/polar";
import { getUserByEmail } from "@/lib/db";

export async function GET() {
  const session = await auth();

  if (!session?.user?.email) {
    return NextResponse.json({ error: "Unauthorized" }, { status: 401 });
  }

  const dbUser = getUserByEmail(session.user.email);
  if (!dbUser?.polar_customer_id) {
    return NextResponse.json({ key: null });
  }

  try {
    const keys = await polar.licenseKeys.list({
      organizationId: process.env.POLAR_ORGANIZATION_ID,
      limit: 100,
    });

    const customerKey = keys.result.items.find(
      (k) => k.customerId === dbUser.polar_customer_id
    );

    return NextResponse.json({ key: customerKey?.key ?? null });
  } catch (e) {
    console.error("Failed to fetch license keys:", e);
    return NextResponse.json({ key: null });
  }
}

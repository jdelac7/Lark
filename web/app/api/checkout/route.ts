import { NextRequest, NextResponse } from "next/server";
import { auth } from "@/lib/auth";
import { polar } from "@/lib/polar";

export async function GET(request: NextRequest) {
  const session = await auth();

  if (!session?.user?.email) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("callbackUrl", request.url);
    return NextResponse.redirect(loginUrl);
  }

  const productId = request.nextUrl.searchParams.get("products");
  if (!productId) {
    return NextResponse.json({ error: "Missing product ID" }, { status: 400 });
  }

  const checkout = await polar.checkouts.create({
    productId,
    customerEmail: session.user.email,
    successUrl: `${process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"}/success?checkout_id={CHECKOUT_ID}`,
  });

  return NextResponse.redirect(checkout.url);
}

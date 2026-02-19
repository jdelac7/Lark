import {
  validateEvent,
  WebhookVerificationError,
} from "@polar-sh/sdk/webhooks";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.text();
  const webhookSecret = process.env.POLAR_WEBHOOK_SECRET;

  if (!webhookSecret) {
    console.error("POLAR_WEBHOOK_SECRET is not set");
    return NextResponse.json(
      { error: "Webhook not configured" },
      { status: 500 }
    );
  }

  let event;
  try {
    event = validateEvent(
      body,
      Object.fromEntries(request.headers),
      webhookSecret
    );
  } catch (e) {
    if (e instanceof WebhookVerificationError) {
      console.error("Webhook verification failed:", e.message);
      return NextResponse.json({ error: "Invalid signature" }, { status: 403 });
    }
    throw e;
  }

  switch (event.type) {
    case "order.created":
      console.log(
        `[Polar] Order created: ${event.data.id}, customer: ${event.data.customerId}`
      );
      break;

    case "subscription.revoked":
      console.log(
        `[Polar] Subscription revoked: ${event.data.id}, customer: ${event.data.customerId}`
      );
      break;

    default:
      console.log(`[Polar] Unhandled event type: ${event.type}`);
  }

  return NextResponse.json({ received: true });
}

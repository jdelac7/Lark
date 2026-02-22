import {
  validateEvent,
  WebhookVerificationError,
} from "@polar-sh/sdk/webhooks";
import { NextRequest, NextResponse } from "next/server";
import { linkPolarCustomer, setSubscriptionStatus } from "@/lib/db";

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
    case "order.created": {
      const customerId = event.data.customerId;
      const customerEmail = event.data.customer?.email;
      console.log(
        `[Polar] Order created: ${event.data.id}, customer: ${customerId}`
      );
      if (customerEmail && customerId) {
        linkPolarCustomer(customerEmail, customerId);
        setSubscriptionStatus(customerId, true);
      }
      break;
    }

    case "subscription.active": {
      const customerId = event.data.customerId;
      console.log(
        `[Polar] Subscription active: ${event.data.id}, customer: ${customerId}`
      );
      if (customerId) {
        setSubscriptionStatus(customerId, true);
      }
      break;
    }

    case "subscription.revoked": {
      const customerId = event.data.customerId;
      console.log(
        `[Polar] Subscription revoked: ${event.data.id}, customer: ${customerId}`
      );
      if (customerId) {
        setSubscriptionStatus(customerId, false);
      }
      break;
    }

    default:
      console.log(`[Polar] Unhandled event type: ${event.type}`);
  }

  return NextResponse.json({ received: true });
}

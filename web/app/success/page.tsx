import { polar } from "@/lib/polar";
import SuccessContent from "./SuccessContent";

interface SuccessPageProps {
  searchParams: Promise<{ checkout_id?: string }>;
}

export default async function SuccessPage({ searchParams }: SuccessPageProps) {
  const { checkout_id } = await searchParams;

  let licenseKey: string | null = null;
  let error: string | null = null;

  if (checkout_id) {
    try {
      const checkout = await polar.checkouts.get({ id: checkout_id });

      if (checkout.customerId) {
        const orgId = process.env.POLAR_ORGANIZATION_ID;
        const keys = await polar.licenseKeys.list({
          organizationId: orgId,
          limit: 100,
        });

        const customerKey = keys.result.items.find(
          (k) => k.customerId === checkout.customerId
        );
        if (customerKey) {
          licenseKey = customerKey.key;
        }
      }
    } catch (e) {
      error =
        "Unable to retrieve your license key. Please check your email for the key, or contact support.";
      console.error("Failed to fetch checkout:", e);
    }
  }

  return <SuccessContent licenseKey={licenseKey} error={error} />;
}

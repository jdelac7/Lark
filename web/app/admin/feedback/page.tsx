import { getAllFeedback } from "@/lib/db";
import Link from "next/link";

export const dynamic = "force-dynamic";

export default function FeedbackPage() {
  const feedback = getAllFeedback();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-lg font-bold text-accent">feedback</h1>
          <p className="text-xs text-text-dim">
            {feedback.length} message{feedback.length !== 1 ? "s" : ""} from
            users
          </p>
        </div>
        <Link
          href="/admin"
          className="rounded border border-border px-3 py-1.5 text-xs text-text-dim transition-colors hover:bg-bg-secondary hover:text-text"
        >
          back to dashboard
        </Link>
      </div>

      <div className="overflow-x-auto rounded border border-border">
        <table className="w-full text-left text-xs">
          <thead className="border-b border-border bg-bg-secondary text-text-dim">
            <tr>
              <th className="px-4 py-2 w-36">date</th>
              <th className="px-4 py-2 w-48">from</th>
              <th className="px-4 py-2">message</th>
            </tr>
          </thead>
          <tbody>
            {feedback.map((f) => (
              <tr
                key={f.id}
                className="border-b border-border/50 transition-colors hover:bg-bg-secondary/50"
              >
                <td className="px-4 py-3 text-text-dim whitespace-nowrap align-top">
                  {f.created_at}
                </td>
                <td className="px-4 py-3 align-top">
                  {f.email ? (
                    <span className="text-cyan">{f.email}</span>
                  ) : (
                    <span className="text-text-dim">anonymous</span>
                  )}
                </td>
                <td className="px-4 py-3 text-text whitespace-pre-wrap">
                  {f.message}
                </td>
              </tr>
            ))}
            {feedback.length === 0 && (
              <tr>
                <td
                  colSpan={3}
                  className="px-4 py-8 text-center text-text-dim"
                >
                  no feedback yet
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

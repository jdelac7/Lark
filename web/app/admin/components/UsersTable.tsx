interface UserRow {
  id: string;
  email: string;
  name: string | null;
  subscribed: number;
  is_admin: number;
  created_at: string;
  ai_cost: number;
}

export default function UsersTable({ users }: { users: UserRow[] }) {
  return (
    <div className="overflow-x-auto rounded border border-border">
      <table className="w-full text-left text-xs">
        <thead className="border-b border-border bg-bg-secondary text-text-dim">
          <tr>
            <th className="px-4 py-2">email</th>
            <th className="px-4 py-2">name</th>
            <th className="px-4 py-2">subscribed</th>
            <th className="px-4 py-2">admin</th>
            <th className="px-4 py-2 text-right">ai cost</th>
            <th className="px-4 py-2">joined</th>
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr
              key={u.id}
              className="border-b border-border/50 transition-colors hover:bg-bg-secondary/50"
            >
              <td className="px-4 py-2 text-text">{u.email}</td>
              <td className="px-4 py-2 text-text-dim">
                {u.name || "---"}
              </td>
              <td className="px-4 py-2">
                {u.subscribed ? (
                  <span className="text-green">yes</span>
                ) : (
                  <span className="text-text-dim">no</span>
                )}
              </td>
              <td className="px-4 py-2">
                {u.is_admin ? (
                  <span className="text-cyan">yes</span>
                ) : (
                  <span className="text-text-dim">no</span>
                )}
              </td>
              <td className="px-4 py-2 text-right font-mono text-yellow">
                ${u.ai_cost.toFixed(4)}
              </td>
              <td className="px-4 py-2 text-text-dim">
                {u.created_at?.slice(0, 10) ?? "---"}
              </td>
            </tr>
          ))}
          {users.length === 0 && (
            <tr>
              <td
                colSpan={6}
                className="px-4 py-8 text-center text-text-dim"
              >
                no users found
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

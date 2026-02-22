interface StatCard {
  label: string;
  value: string | number;
  sub?: string;
  color?: string;
}

export default function AdminStats({ stats }: { stats: StatCard[] }) {
  return (
    <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
      {stats.map((stat) => (
        <div
          key={stat.label}
          className="rounded border border-border bg-bg-card p-4"
        >
          <div className="mb-1 text-xs text-text-dim">{stat.label}</div>
          <div
            className={`text-xl font-bold ${stat.color ?? "text-accent"}`}
          >
            {stat.value}
          </div>
          {stat.sub && (
            <div className="mt-1 text-xs text-text-dim">{stat.sub}</div>
          )}
        </div>
      ))}
    </div>
  );
}

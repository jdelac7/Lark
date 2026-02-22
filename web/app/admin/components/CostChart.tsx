"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";

interface DailyData {
  date: string;
  premade_cost: number;
  custom_cost: number;
}

export default function CostChart({ data }: { data: DailyData[] }) {
  if (data.length === 0) {
    return (
      <div className="flex h-64 items-center justify-center rounded border border-border bg-bg-card text-sm text-text-dim">
        no cost data yet
      </div>
    );
  }

  return (
    <div className="rounded border border-border bg-bg-card p-4">
      <ResponsiveContainer width="100%" height={300}>
        <AreaChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#1e1e1e" />
          <XAxis
            dataKey="date"
            stroke="#555"
            tick={{ fontSize: 10, fill: "#555" }}
            tickFormatter={(v: string) => v.slice(5)}
          />
          <YAxis
            stroke="#555"
            tick={{ fontSize: 10, fill: "#555" }}
            tickFormatter={(v: number) => `$${v.toFixed(2)}`}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: "#141414",
              border: "1px solid #1e1e1e",
              borderRadius: 4,
              fontSize: 12,
              fontFamily: "JetBrains Mono, monospace",
              color: "#b0b0b0",
            }}
            formatter={(value, name) => [
              `$${(Number(value) || 0).toFixed(4)}`,
              name === "premade_cost" ? "premade" : "custom",
            ]}
            labelFormatter={(label) => `date: ${label}`}
          />
          <Legend
            wrapperStyle={{ fontSize: 11, fontFamily: "JetBrains Mono" }}
            formatter={(value: string) =>
              value === "premade_cost" ? "premade" : "custom"
            }
          />
          <Area
            type="monotone"
            dataKey="premade_cost"
            stackId="1"
            stroke="#00ff41"
            fill="#00ff41"
            fillOpacity={0.15}
          />
          <Area
            type="monotone"
            dataKey="custom_cost"
            stackId="1"
            stroke="#a78bfa"
            fill="#a78bfa"
            fillOpacity={0.15}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}

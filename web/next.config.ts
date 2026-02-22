import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  rewrites: async () => [
    {
      source: "/game-api/:path*",
      destination: `${process.env.GAME_SERVER_INTERNAL_URL || "http://localhost:9292"}/api/v1/:path*`,
    },
  ],
};

export default nextConfig;

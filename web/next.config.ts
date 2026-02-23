import type { NextConfig } from "next";

const gameServer =
  process.env.GAME_SERVER_INTERNAL_URL || "http://localhost:9292";

const nextConfig: NextConfig = {
  output: "standalone",
  rewrites: async () => ({
    // beforeFiles rewrites are checked before pages/routes — keep empty
    beforeFiles: [],
    // afterFiles rewrites are checked AFTER filesystem routes (App Router)
    // so our route handlers at /game-api/scenarios/start/stream and
    // /game-api/game/input/stream take priority over these catch-alls.
    afterFiles: [
      {
        source: "/game-api/:path*",
        destination: `${gameServer}/api/v1/:path*`,
      },
      {
        source: "/api/v1/:path*",
        destination: `${gameServer}/api/v1/:path*`,
      },
    ],
    fallback: [],
  }),
};

export default nextConfig;

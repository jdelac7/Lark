# ── Stage 1: Build Go server ──────────────────────────────────────────
FROM golang:1.24-bookworm AS go-build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY api/ api/
COPY server/ server/
RUN CGO_ENABLED=0 go build -o /lark-server ./server

# ── Stage 2: Build Next.js web app ───────────────────────────────────
FROM node:20-bookworm-slim AS web-build

WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./

# Provide build-time defaults so next build succeeds
ENV AUTH_SECRET=build-placeholder

RUN npm run build

# ── Stage 3: Production image ────────────────────────────────────────
FROM node:20-bookworm-slim AS production

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates bash \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Go server binary
COPY --from=go-build /lark-server ./lark-server

# Next.js standalone output
COPY --from=web-build /src/web/.next/standalone ./web/
COPY --from=web-build /src/web/.next/static ./web/.next/static
COPY --from=web-build /src/web/public ./web/public

# Data directories (mounted as volumes in production)
RUN mkdir -p /app/data /app/web/data

# Entrypoint script
COPY <<'ENTRYPOINT' /app/start.sh
#!/bin/bash
set -e

# Symlink shared data directory so both services see the same DBs
ln -sf /app/data/lark.db /app/web/data/lark.db 2>/dev/null || true

# Start Go server in background
cd /app
PORT=${GAME_SERVER_PORT:-9292} ./lark-server &
GO_PID=$!

# Start Next.js
cd /app/web
HOSTNAME=0.0.0.0 \
PORT=${WEB_PORT:-3000} \
COST_DB_PATH=/app/cost.db \
GAME_SERVER_INTERNAL_URL=http://localhost:9292 \
node server.js &
WEB_PID=$!

# Wait for either to exit
wait -n $GO_PID $WEB_PID
EXIT_CODE=$?

# If one dies, kill the other
kill $GO_PID $WEB_PID 2>/dev/null || true
exit $EXIT_CODE
ENTRYPOINT
RUN chmod +x /app/start.sh

# Volumes for persistent data
VOLUME ["/app/data"]

EXPOSE 3000 9292

ENV NODE_ENV=production

CMD ["/app/start.sh"]

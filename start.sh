#!/bin/bash
set -e

# Symlink shared data directory so both services see the same DBs
ln -sf /app/data/lark.db /app/web/data/lark.db 2>/dev/null || true

# Force PORT for each service to avoid PaaS conflicts
export GO_PORT="${GAME_SERVER_PORT:-9292}"
export WEB_PORT="${WEB_PORT:-3000}"

echo "Starting Go server on :${GO_PORT}"
echo "Starting Next.js on :${WEB_PORT}"

# Start Go server in background
cd /app
PORT=${GO_PORT} ./lark-server &
GO_PID=$!

# Start Next.js
cd /app/web
HOSTNAME=0.0.0.0 \
PORT=${WEB_PORT} \
COST_DB_PATH=/app/cost.db \
GAME_SERVER_INTERNAL_URL=http://localhost:${GO_PORT} \
node server.js &
WEB_PID=$!

# Wait for either to exit
wait -n $GO_PID $WEB_PID
EXIT_CODE=$?

# If one dies, kill the other
kill $GO_PID $WEB_PID 2>/dev/null || true
exit $EXIT_CODE

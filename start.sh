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

# Write .env for Next.js standalone — it doesn't inherit Docker env vars reliably
cd /app/web
cat > .env <<EOF
AUTH_URL=${AUTH_URL:-${NEXTAUTH_URL:-https://lark.black}}
AUTH_TRUST_HOST=true
AUTH_SECRET=${AUTH_SECRET}
AUTH_GOOGLE_ID=${AUTH_GOOGLE_ID}
AUTH_GOOGLE_SECRET=${AUTH_GOOGLE_SECRET}
POLAR_ACCESS_TOKEN=${POLAR_ACCESS_TOKEN}
POLAR_WEBHOOK_SECRET=${POLAR_WEBHOOK_SECRET}
POLAR_ORGANIZATION_ID=${POLAR_ORGANIZATION_ID}
NEXT_PUBLIC_POLAR_PRODUCT_ID=${NEXT_PUBLIC_POLAR_PRODUCT_ID}
NEXT_PUBLIC_SITE_URL=${NEXT_PUBLIC_SITE_URL:-https://lark.black}
COST_DB_PATH=/app/cost.db
GAME_SERVER_INTERNAL_URL=http://localhost:${GO_PORT}
EOF

echo "Wrote .env with $(wc -l < .env) vars"

HOSTNAME=0.0.0.0 PORT=${WEB_PORT} node server.js &
WEB_PID=$!

# Wait for either to exit
wait -n $GO_PID $WEB_PID
EXIT_CODE=$?

# If one dies, kill the other
kill $GO_PID $WEB_PID 2>/dev/null || true
exit $EXIT_CODE

#!/bin/bash

set -e

echo "🧹 Resetting database..."

# load env safely
set -a
source .env
set +a

if [ -z "$DATABASE_URL" ]; then
  echo "❌ DATABASE_URL not found in .env"
  exit 1
fi

# drop and recreate schema
psql "$DATABASE_URL" <<EOF
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
EOF

echo "📦 Running migrations..."

make migrate-up

echo "✅ Database reset complete!"
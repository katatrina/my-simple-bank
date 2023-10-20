#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" --verbose up # execute the "migrate" binary

echo "start the app"
exec "$@"
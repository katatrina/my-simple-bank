#!/bin/sh

set -e

echo "run db migrations"
source /app/app.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" --verbose up # execute the "migrate" binary

echo "start the app"
exec "$@"
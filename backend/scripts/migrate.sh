#!/bin/bash
# Run all pending SQL migrations in order
set -euo pipefail

MIGRATIONS_DIR="$(dirname "$0")/../internal/database/migrations"

if [ -z "${DB_HOST:-}" ] || [ -z "${DB_USER:-}" ] || [ -z "${DB_PASSWORD:-}" ] || [ -z "${DB_NAME:-}" ]; then
    echo "Error: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME must be set"
    exit 1
fi

echo "Running migrations against $DB_HOST/$DB_NAME..."

for file in $(ls "$MIGRATIONS_DIR"/*.sql | sort); do
    echo "Applying: $(basename $file)..."
    mysql -h "$DB_HOST" -P "${DB_PORT:-3306}" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$file"
    echo "  Done."
done

echo "All migrations applied successfully."

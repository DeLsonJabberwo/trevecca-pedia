#!/bin/bash
# setup-db.sh - Apply database schema to fly.io Postgres
# This script only applies schema files 01 and 02 (structure), skipping 03 (seed data)
# You must handle postgres connection and secrets configuration separately

set -e

DB_APP_NAME="${1:-trevecca-pedia-db}"
DB_NAME="${2:-trevecca_pedia_wiki}"  # Default to the database created by fly postgres attach

echo "Applying schema to database app: $DB_APP_NAME"
echo "Target database: $DB_NAME"
echo ""

# Check if fly CLI is installed
if ! command -v fly &> /dev/null; then
    echo "Error: fly CLI is not installed"
    echo "Install it from: https://fly.io/docs/hands-on/install-flyctl/"
    exit 1
fi

# Check if we're logged in
if ! fly auth whoami &> /dev/null; then
    echo "Error: Not logged into fly.io"
    echo "Run: fly auth login"
    exit 1
fi

# Check if database app exists
if ! fly status --app "$DB_APP_NAME" &> /dev/null; then
    echo "Error: Database app '$DB_APP_NAME' not found"
    echo "Create it first: fly postgres create --name $DB_APP_NAME"
    exit 1
fi

# Apply schema files 01 and 02 only (not 03 which contains seed data)
echo "Applying schema files..."

for file in init/01-schema.sql init/02-schema.sql; do
    if [ -f "$file" ]; then
		echo "Applying $file to database '$DB_NAME'..."
        fly postgres connect --app "$DB_APP_NAME" --database "$DB_NAME" < "$file"
        echo "✓ Applied $file"
    else
        echo "Warning: $file not found"
    fi
done

echo ""
echo "========================================="
echo "Schema applied successfully!"
echo "========================================="
echo ""
echo "Note: Seed data (init/03-schema.sql) was NOT applied."
echo ""
echo "Next steps for connecting your wiki service:"
echo "1. Get connection info: fly postgres connect --app $DB_APP_NAME --command \"\\conninfo\""
echo "2. Set secrets manually or use: fly postgres attach $DB_APP_NAME --app <wiki-app-name>"
echo ""

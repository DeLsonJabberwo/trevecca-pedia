#!/bin/bash
# setup-db.sh - Initialize database schema on fly.io Postgres
# This script applies schema files 01 and 02 (structure), skipping 03 (seed data)

set -e

DB_APP_NAME="${1:-trevecca-pedia-db}"

echo "Setting up database for: $DB_APP_NAME"
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

# Get the database connection URL
echo "Getting database connection URL..."
DB_URL=$(fly postgres connect --app "$DB_APP_NAME" --command "\\conninfo" 2>/dev/null || echo "")

if [ -z "$DB_URL" ]; then
    echo "Error: Could not get database connection info"
    echo "Make sure your Postgres app is running: fly status --app $DB_APP_NAME"
    exit 1
fi

echo "Connection info retrieved"
echo ""

# Apply schema files 01 and 02 only (not 03 which contains seed data)
echo "Applying schema..."

for file in init/01-schema.sql init/02-schema.sql; do
    if [ -f "$file" ]; then
        echo "Applying $file..."
        fly postgres connect --app "$DB_APP_NAME" < "$file"
        echo "✓ Applied $file"
    else
        echo "Warning: $file not found"
    fi
done

echo ""
echo "========================================="
echo "Database schema applied successfully!"
echo "========================================="
echo ""
echo "Note: Seed data (init/03-schema.sql) was NOT applied."
echo "The database is ready for use with an empty schema."
echo ""
echo "Next steps:"
echo "1. Get the database connection string:"
echo "   fly postgres connect --app $DB_APP_NAME --command \"\\conninfo\""
echo ""
echo "2. Set the secrets for your wiki service:"
echo "   cd ../wiki"
echo "   fly secrets set WIKI_DB_HOST=<host> --app trevecca-pedia-wiki"
echo "   fly secrets set WIKI_DB_PORT=<port> --app trevecca-pedia-wiki"
echo "   fly secrets set WIKI_DB_PASSWORD=<password> --app trevecca-pedia-wiki"

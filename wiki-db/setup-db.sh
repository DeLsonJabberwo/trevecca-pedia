#!/bin/bash
# setup-db.sh - Apply database schema to fly.io Postgres
# This script only applies schema files 01 and 02 (structure), skipping 03 (seed data)
# NOTE: You must run 'fly postgres attach' before this script to create the database

set -e

DB_APP_NAME="${1:-trevecca-pedia-db}"
DB_NAME="${2:-trevecca_pedia_wiki}"

echo "========================================="
echo "Database Schema Setup"
echo "Database app: $DB_APP_NAME"
echo "Target database: $DB_NAME"
echo "========================================="
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

echo "Applying schema to database '$DB_NAME'..."
echo ""

# Apply schema files 01 and 02 only (not 03 which contains seed data)
echo "Applying schema files..."

for file in init/01-schema.sql init/02-schema.sql; do
    if [ -f "$file" ]; then
        echo "  Applying $file..."
        # Use printf to ensure proper newline handling and add \q to exit psql
        printf '%s\n\\q\n' "$(cat "$file")" | fly postgres connect --app "$DB_APP_NAME" --database "$DB_NAME"
        echo "  ✓ Applied $file"
    else
        echo "  Warning: $file not found, skipping"
    fi
done

echo ""
echo "========================================="
echo "Schema applied successfully!"
echo "========================================="
echo ""
echo "Note: Seed data (init/03-schema.sql) was NOT applied."
echo "      To apply seed data, run:"
echo "      printf '%s\\n\\\\q\\n' \"\$(cat init/03-schema.sql)\" | fly postgres connect --app $DB_APP_NAME --database $DB_NAME"
echo ""

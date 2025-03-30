#!/bin/bash
set -e

echo "Checking if schema.sql is in sync with migrations..."

# Create a temp file for comparison
TEMP_SCHEMA=$(mktemp)
SCHEMA_FILE="db/schema.sql"

# Run the same generation process as generate_schema.sh but output to temp file
# Create a header for the schema file
cat >"$TEMP_SCHEMA" <<EOL
-- Database schema for Gopher Tower
-- AUTOMATICALLY GENERATED FROM MIGRATIONS
-- DO NOT EDIT DIRECTLY
--
-- Generated: $(date)
--
-- This schema is used by sqlc to generate Go code
-- for database operations

EOL

# Use sqlite to apply migrations and dump schema
sqlite3 ":memory:" <<EOL
-- Apply migrations in order
$(find ./internal/db/migrations -name "*.up.sql" | sort | xargs cat)

-- Output schema as SQL statements
.output ${TEMP_SCHEMA}.tmp
.schema
.quit
EOL

# Process the schema output to make it more readable
cat "${TEMP_SCHEMA}".tmp | grep -v "CREATE TABLE sqlite_" >>"$TEMP_SCHEMA"
rm "${TEMP_SCHEMA}".tmp

# Strip date lines from both files for comparison
TEMP_NO_DATE=$(mktemp)
SCHEMA_NO_DATE=$(mktemp)

grep -v "Generated:" "$TEMP_SCHEMA" >"$TEMP_NO_DATE"
grep -v "Generated:" $SCHEMA_FILE >"$SCHEMA_NO_DATE"

# Compare files (ignoring the date line)
if diff -q "$TEMP_NO_DATE" "$SCHEMA_NO_DATE" >/dev/null; then
  echo "✅ Schema is in sync with migrations"
  exit 0
else
  echo "❌ ERROR: schema.sql is out of sync with migrations"
  echo "Please run 'task schema:generate' to update it"
  echo ""
  echo "Diff:"
  diff -u "$SCHEMA_NO_DATE" "$TEMP_NO_DATE" || true

  # Clean up
  rm -f "$TEMP_SCHEMA" "$TEMP_NO_DATE" "$SCHEMA_NO_DATE"
  exit 1
fi

# Clean up
rm -f "$TEMP_SCHEMA" "$TEMP_NO_DATE" "$SCHEMA_NO_DATE"

#!/bin/bash
set -e

echo "Generating canonical schema.sql from migrations..."

# Create a temporary database
TEMP_DB=":memory:"

# Path to migrations
MIGRATIONS_DIR="internal/db/migrations"

# Output file
OUTPUT_FILE="db/schema.sql"

# Create a header for the schema file
cat >$OUTPUT_FILE <<EOL
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
sqlite3 $TEMP_DB <<EOL
-- Apply migrations in order
$(find ./$MIGRATIONS_DIR -name "*.up.sql" | sort | xargs cat)

-- Output schema as SQL statements
.output ${OUTPUT_FILE}.tmp
.schema
.quit
EOL

# Process the schema output to make it more readable
cat ${OUTPUT_FILE}.tmp | grep -v "CREATE TABLE sqlite_" >>$OUTPUT_FILE
rm ${OUTPUT_FILE}.tmp

echo "Schema generation complete: $OUTPUT_FILE"

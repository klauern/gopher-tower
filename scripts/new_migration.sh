#!/bin/bash
set -e

if [ $# -lt 1 ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi

# Migration name from argument
MIGRATION_NAME=$1
TIMESTAMP=$(date +%Y%m%d%H%M%S)
MIGRATION_FILE="internal/db/migrations/$(printf '%06d' "$TIMESTAMP")_${MIGRATION_NAME}"

# Create migration files
echo "Creating migration: $MIGRATION_NAME"
touch "${MIGRATION_FILE}.up.sql"
touch "${MIGRATION_FILE}.down.sql"

echo "Created migration files:"
echo "  ${MIGRATION_FILE}.up.sql"
echo "  ${MIGRATION_FILE}.down.sql"

cat <<"EOF"

⚠️  IMPORTANT WORKFLOW STEPS ⚠️

1. Edit the migration files to define your changes
2. Run 'task schema:generate' to update schema.sql
3. Run 'task schema:check' to verify everything is in sync
4. Include both migration files AND updated schema.sql in your commit

NOTE: The CI pipeline will fail if schema.sql is out of sync with migrations!
EOF

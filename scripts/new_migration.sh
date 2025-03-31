#!/bin/bash
set -e

# Check if golang-migrate is installed
if ! command -v migrate &>/dev/null; then
  echo "Error: golang-migrate is not installed"
  echo "Install it with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  exit 1
fi

prompt_migration_name() {
  local name=""
  while [ -z "$name" ]; do
    read -p "Enter migration name (e.g., add_users_table): " name
    # Remove spaces and special characters, convert to snake_case
    name=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr ' ' '_' | sed 's/[^a-z0-9_]//g')
    if [ -z "$name" ]; then
      echo "Error: Migration name cannot be empty"
    fi
  done
  echo "$name"
}

# Get migration name from argument or prompt
MIGRATION_NAME=""
if [ $# -ge 1 ]; then
  MIGRATION_NAME=$1
else
  MIGRATION_NAME=$(prompt_migration_name)
fi

# Validate migration name
if [[ ! $MIGRATION_NAME =~ ^[a-z0-9_]+$ ]]; then
  echo "Error: Migration name can only contain lowercase letters, numbers, and underscores"
  exit 1
fi

# Create migration using golang-migrate
echo "Creating migration: $MIGRATION_NAME"
migrate create -ext sql -dir internal/db/migrations -seq "$MIGRATION_NAME"

# Find the newly created migration files
MIGRATION_FILES=($(ls -t internal/db/migrations/*"${MIGRATION_NAME}".*.sql | head -n 2))
UP_MIGRATION=${MIGRATION_FILES[0]}
DOWN_MIGRATION=${MIGRATION_FILES[1]}

echo "Created migration files:"
echo "  $UP_MIGRATION"
echo "  $DOWN_MIGRATION"

# Prompt for migration content
read -p "Would you like to edit the migration files now? [Y/n] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
  # Try to use the user's preferred editor
  EDITOR=${EDITOR:-$(which vim || which nano || which vi)}
  if [ -n "$EDITOR" ]; then
    $EDITOR "$UP_MIGRATION"
    echo "Now edit the down migration..."
    $EDITOR "$DOWN_MIGRATION"
  else
    echo "No text editor found. Please edit the files manually:"
    echo "  $UP_MIGRATION"
    echo "  $DOWN_MIGRATION"
  fi
fi

cat <<"EOF"

⚠️  IMPORTANT WORKFLOW STEPS ⚠️

1. Edit the migration files to define your changes
2. Run 'task schema:generate' to update schema.sql
3. Run 'task schema:check' to verify everything is in sync
4. Include both migration files AND updated schema.sql in your commit

NOTE: The CI pipeline will fail if schema.sql is out of sync with migrations!
EOF

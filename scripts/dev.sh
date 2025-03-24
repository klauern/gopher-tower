#!/usr/bin/env bash

# Exit on error
set -e

# Cleanup function to ensure both processes are stopped
cleanup() {
  echo "Shutting down development servers..."

  # Kill the Go server if it's running
  if [ -n "$GO_PID" ]; then
    kill $GO_PID 2>/dev/null || true
  fi

  # Kill any remaining processes in our process group
  pkill -P $$ 2>/dev/null || true

  exit 0
}

# Set up trap for cleanup
trap cleanup EXIT INT TERM

# Start the Go server
echo "Starting Go server..."
go run ./cmd/main.go &
GO_PID=$!

# Start the frontend dev server
echo "Starting frontend dev server..."
cd frontend && exec npm run dev

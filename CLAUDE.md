# Gopher Tower Project Guide

## Build & Test Commands
- **Frontend Dev**: `task frontend:dev` or `cd frontend && bun dev`
- **Backend Dev**: `task go:run` or `go run ./cmd/main.go`
- **Full Dev Environment**: `task dev` (runs both frontend and backend)
- **Lint**: `task lint` (or `task go:lint` / `task frontend:lint`)
- **Test**: `task test` (or `task go:test` / `task frontend:test`)
- **Build**: `task build` (or `task go:build` / `task frontend:build`)
- **Single Go test**: `go test ./path/to/package -v -run TestName`
- **Embedded Build**: `task embedded` (builds single binary with frontend embedded)
- **Run Embedded**: `task run:embedded` (builds and runs embedded application)

## Code Style Guidelines
- **Go**: Follow standard Go conventions (gofmt)
  - Error handling: Always check errors and provide context
  - Types: Use strong typing with meaningful names
  - Imports: Group standard lib, external, then internal
  - Use embed package for including frontend assets

- **Frontend (TypeScript/React)**:
  - Use TypeScript interfaces/types for all components
  - Prefer functional components with hooks
  - Follow ESLint rules from next/core-web-vitals
  - Use Tailwind for styling
  - Client components: mark with 'use client' directive
  - Use relative URLs for API endpoints

This project combines a Go backend (SSE server) with a Next.js frontend using App Router and TypeScript. The application can run in separate processes during development or as a single embedded binary for production.

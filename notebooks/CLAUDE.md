# Gopher Tower Project Guide

## Build & Test Commands

- **Dev Environment**: `task dev` (runs both frontend and backend), `task frontend:dev`, `task go:dev`
- **Build**: `task build` (or `task go:build` / `task frontend:build`); use `task embedded` for single binary
- **Lint**: `task lint` (or `task go:lint` / `task frontend:lint`)
- **Format**: `task fmt` (or `task go:fmt` / `task frontend:fmt`)
- **Test**: `task test` (or `task go:test` / `task frontend:test`)
- **Single Frontend Test**: `cd frontend && npm run test -- -t "test name pattern"`
- **Single Go Test**: `go test ./path/to/package -v -run TestName`
- **Test Coverage**: `task test:coverage` (or `task go:test:coverage` / `task frontend:test:coverage`)

## Development Environment

### DevContainer Setup

- Project includes VS Code DevContainer configuration for consistent development
- Key extensions pre-configured:
  - Go tools and language server
  - ESLint and Prettier for TypeScript/JavaScript
  - Tailwind CSS support
  - TypeScript and React tools
  - GitHub Copilot and GitLens
- Features enabled:
  - Git and GitHub CLI integration
  - Latest Go version
- Automatic port forwarding: 3000 (frontend) and 8080 (backend)
- Dependencies installed automatically via `task deps`
- Git configuration mounted from host system
- Workspace mounted with delegated consistency for better performance
- Lefthook installed automatically for Git hooks
- Cursor server workspace storage configured

### Git Hooks with Lefthook

- Lefthook is automatically installed in the development container
- Provides pre-commit and pre-push hooks for:
  - Code formatting checks
  - Linting
  - Type checking
  - Test execution
- Configure hooks in `lefthook.yml` at project root

## Architecture Overview

- Go backend + Next.js 15 frontend (using App Router) with Tailwind and TypeScript
- Frontend is built as a fully static site (`output: 'export'`), embedded and served by the Go backend
- React components utilize the Next.js App Router model (Server Components by default, Client Components with 'use client')
- Plugin framework allows running various job types (CLI tools, Python scripts)
- SQLite database with migration support using golang-migrate

## Code Style Guidelines

### Go

- Follow standard Go conventions with `gofmt`
- Error handling: Always check errors and provide context with `fmt.Errorf("context: %w", err)`
- Types: Strong typing with meaningful names; use interfaces for abstractions
- Imports: Group standard lib, external, then internal packages
- Tests: Located next to implementation files with detailed table-driven tests

### Frontend (TypeScript/React)

- Use TypeScript for all code with proper interfaces/types
- Components: Follow Next.js App Router conventions (Server and Client Components).
- Styling: Tailwind CSS following the component patterns in `frontend/components/ui`
- Tests: Using Vitest + React Testing Library, typically located in `__tests__` directories alongside components (e.g., `frontend/components/__tests__`) or colocated with `.test.tsx` / `.test.ts` extensions.
- API endpoints: Use relative URLs (e.g., `/api/jobs`). The Next.js development server proxies `/api/*` requests to the Go backend (see `frontend/next.config.ts`). In production, the Go server handles these directly.

## Plugin Framework

- Plugin system supports multiple job execution engines (CLI tools, Python scripts)
- Plugins implement the Plugin interface with Name(), Execute(), Validate(), etc.
- New plugins should follow security best practices, including validation and resource control
- See `docs/PLUGIN_FRAMEWORK.md` for detailed implementation information

## Database Migrations

- Uses `golang-migrate` for database schema versioning and migrations
- Migration files are stored in `internal/db/migrations`
- Migrations are automatically run on application startup
- Migration state is tracked in SQLite database
- Each migration has an up and down file for forward and rollback changes
- Migration files follow the format: `{version}_{description}.{up|down}.sql`

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

- Go backend + Next.js 15 frontend with Tailwind and TypeScript
- Frontend is built as a fully static site, embedded and served by Go backend
- All React components MUST use 'use client' directive (client components only)
- Plugin framework allows running various job types (CLI tools, Python scripts)

## Code Style Guidelines

### Go
- Follow standard Go conventions with `gofmt`
- Error handling: Always check errors and provide context with `fmt.Errorf("context: %w", err)`
- Types: Strong typing with meaningful names; use interfaces for abstractions
- Imports: Group standard lib, external, then internal packages
- Tests: Located next to implementation files with detailed table-driven tests

### Frontend (TypeScript/React)
- Use TypeScript for all code with proper interfaces/types
- Components: Client-side only with 'use client' directive
- Styling: Tailwind CSS following the component patterns in `frontend/components/ui`
- Tests: Using Vitest + React Testing Library in `__tests__` directories or next to components
- API endpoints: Use relative URLs, determined at runtime from window.location

## Plugin Framework
- Plugin system supports multiple job execution engines (CLI tools, Python scripts)
- Plugins implement the Plugin interface with Name(), Execute(), Validate(), etc.
- New plugins should follow security best practices, including validation and resource control
- See `docs/PLUGIN_FRAMEWORK.md` for detailed implementation information

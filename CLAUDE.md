# Gopher Tower Project Guide

## Build & Test Commands

- **Frontend Dev**: `task frontend:dev` or `cd frontend && bun dev`
- **Backend Dev**: `task go:run` or `go run ./cmd/main.go`
- **Full Dev Environment**: `task dev` (runs both frontend and backend)
- **Lint**: `task lint` (or `task go:lint` / `task frontend:lint`)
- **Test**: `task test` (or `task go:test` / `task frontend:test`)
- **Test Watch Mode**: `task test:watch` (or `task frontend:test:watch`)
- **Test with Coverage**: `task test:coverage` (or `task frontend:test:coverage`)
- **Test with UI**: `task frontend:test:ui`
- **Build**: `task build` (or `task go:build` / `task frontend:build`)
- **Single Go test**: `go test ./path/to/package -v -run TestName`
- **Embedded Build**: `task embedded` (builds single binary with frontend embedded)
- **Run Embedded**: `task run:embedded` (builds and runs embedded application)

## Architecture Overview

This project uses a static site architecture where:

- Frontend is built as a fully static Next.js application
- All components MUST be client components (marked with 'use client')
- No server-side rendering or server components are used
- Frontend is embedded and served by the Go backend
- API calls use relative URLs that are resolved against the running server

## Next.js 15 Important Notes

- **Route Parameters**: In Next.js 15, route parameters (`params`) are now Promises and must be unwrapped using `React.use()`:

  ```typescript
  // Correct way to handle route params in Next.js 15
  import { use } from 'react';

  export default function Page({ params }: { params: Promise<{ id: string }> }) {
    const resolvedParams = use(params);
    // Now you can safely use resolvedParams.id
  }
  ```

- Direct access to `params` properties (e.g., `params.id`) will trigger warnings and will be deprecated in future versions

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
  - ALL components must be client components with 'use client' directive
  - Base URLs for API endpoints are determined at runtime from window.location
  - Environment variables (like NEXT_PUBLIC_API_BASE_URL) can override the base URL

## Testing

- **Frontend Testing**: Uses Vitest with React Testing Library
  - Test files should be placed in `app/__tests__` directory or next to components
  - Test files should be named with `.test.tsx` or `.test.ts` extension
  - See `frontend/TESTING.md` for detailed documentation and examples
  - Run tests with `bun run test` or `task frontend:test`

- **Backend Testing**: Uses Go's standard testing package
  - Tests are located next to the code they test
  - Run tests with `go test ./...` or `task go:test`

This project combines a Go backend (SSE server) with a Next.js frontend using App Router and TypeScript. The frontend is built as a static site and embedded into the Go binary for production deployment. During development, they can run as separate processes for faster iteration.

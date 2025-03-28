# Gopher Tower

This is a web application combining a Go backend with a Next.js 15 frontend, featuring server-sent events (SSE) capabilities.  Developed almost exclusively using [vibe coding](https://en.wikipedia.org/wiki/Vibe_coding) (seriously).  I just wanted to create a NextJS frontend that could be embedded as part of a Golang backend in one static binary.  The rest of this README is Claude-generated.

## 🚀 Features

- **Frontend (Next.js 15)**
  - TypeScript support
  - Tailwind CSS for styling
  - Pages Router architecture (I refactored from App Router since it's client-only)
  - Client-side components
  - Turbopack for fast development
  - ESLint configuration
  - Comprehensive testing setup with Vitest and React Testing Library

- **Backend (Go)**
  - Server-sent events (SSE) implementation
  - Static file serving
  - Production-ready embedded frontend
  - Standard Go practices and conventions

## 📋 Prerequisites

- Node.js 22.x
- Go 1.x
- [Task](https://taskfile.dev/) - Task runner for development workflows

## 🛠️ Setup and Installation

1. Clone the repository:

   ```bash
   git clone github.com/klauern/gopher-tower.git
   cd gopher-tower
   ```

2. Install all dependencies:

   ```bash
   task deps
   ```

That's it! You're ready to start development.

## 🔧 Development

### Running the Application

- **Full Development Environment**:

  ```bash
  task dev
  ```

  This will start both frontend and backend servers.

- **Frontend Only**:

  ```bash
  task frontend:dev
  # or
  cd frontend && npm run dev
  ```

- **Backend Only**:

  ```bash
  task go:run
  # or
  go run ./cmd/main.go
  ```

### Testing

- **Run All Tests**:

  ```bash
  task test
  ```

- **Frontend Tests**:

  ```bash
  task frontend:test          # Run tests
  task frontend:test:watch    # Watch mode
  task frontend:test:coverage # With coverage
  task frontend:test:ui       # With UI
  ```

- **Backend Tests**:

  ```bash
  task go:test
  ```

- **Single Go Test**:

  ```bash
  go test ./path/to/package -v -run TestName
  ```

### Linting

```bash
task lint              # Lint everything
task go:lint          # Lint Go code
task frontend:lint    # Lint frontend code
```

## 🏗️ Building

- **Complete Build**:

  ```bash
  task build
  ```

- **Embedded Build** (single binary with frontend):

  ```bash
  task embedded
  ```

- **Run Embedded Application**:

  ```bash
  task run:embedded
  ```

## 🏗️ Project Structure

```
.
├── frontend/               # Next.js 15 frontend application
│   ├── app/               # App router pages and components
│   ├── __tests__/        # Frontend tests
│   └── TESTING.md        # Frontend testing documentation
├── cmd/                   # Go application entrypoints
├── Taskfile.yml          # Task runner configuration
└── README.md             # This file
```

## 🧪 Testing Guidelines

### Frontend Testing

- Tests are located in `frontend/app/__tests__` or alongside components
- Use `.test.tsx` or `.test.ts` extension for test files
- Implements Vitest with React Testing Library
- Refer to `frontend/TESTING.md` for detailed testing documentation

### Backend Testing

- Uses Go's standard testing package
- Tests are co-located with the code they test
- Follow standard Go testing conventions

## 🔍 Code Style Guidelines

### Go

- Follow standard Go conventions (gofmt)
- Proper error handling with context
- Strong typing with meaningful names
- Organized imports (standard lib, external, internal)
- Use embed package for frontend assets

### Frontend

- TypeScript interfaces/types for all components
- Functional components with hooks
- ESLint rules from next/core-web-vitals
- Tailwind CSS for styling
- All components must use 'use client' directive
- Dynamic API endpoint resolution based on runtime

## 📝 Important Notes

- The frontend is built as a static site
- All components MUST be client components
- No server-side rendering or server components
- Frontend is embedded and served by the Go backend
- API calls use relative URLs resolved against the running server

## 📄 License

[Add your license information here]

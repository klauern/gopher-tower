# Frontend build stage
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Install frontend dependencies
COPY frontend/package*.json ./
RUN npm install

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Go build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend files to be embedded
COPY --from=frontend-builder /app/frontend/out ./internal/static/frontend/

# Build the application with optimizations for multiple platforms
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
  -ldflags="-w -s" \
  -tags embed \
  -o /app/server \
  ./cmd/main.go

# Final stage - using distroless which is minimal but includes certificates
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy the binary from builder
COPY --from=backend-builder /app/server .

EXPOSE 8080

USER nonroot:nonroot

CMD ["./server"]

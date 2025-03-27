FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm

# Install Node.js 22.x
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
  && apt-get update \
  && apt-get install -y nodejs \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Install task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Install common development tools
RUN apt-get update && apt-get install -y \
  git \
  curl \
  wget \
  vim \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /workspace

# Install Go tools
RUN go install golang.org/x/tools/gopls@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
  && go install github.com/golang/mock/mockgen@latest

# Verify installations
RUN go version && \
  node --version && \
  npm --version && \
  task --version

FROM mcr.microsoft.com/devcontainers/go:1.24

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# # Install Node.js 22.x
# # hadolint ignore=DL3008
# RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
#   && apt-get update \
#   && apt-get install -y nodejs --no-install-recommends \
#   && apt-get clean \
#   && rm -rf /var/lib/apt/lists/*

# Install task
# RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Install common development tools
# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
  git \
  curl \
  wget \
  vim \
  sqlite3 \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /workspace

# # Install Go tools
# RUN go install golang.org/x/tools/gopls@latest \
#   && go install github.com/go-delve/delve/cmd/dlv@latest \
#   && go install github.com/golang/mock/mockgen@latest

# Verify installations
RUN go version

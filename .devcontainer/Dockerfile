FROM mcr.microsoft.com/devcontainers/go:1.24

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install task
# RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
COPY scripts/install_git_delta.sh /tmp/install_git_delta.sh
RUN /tmp/install_git_delta.sh

RUN curl --proto '=https' --tlsv1.2 -LsSf https://setup.atuin.sh | sh

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

# Verify installations
RUN go version

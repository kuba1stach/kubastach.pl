# Devcontainer Dockerfile for Go backend + frontend with OpenAPI tools
FROM ubuntu:22.04

# Install system dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    make \
    unzip \
    wget \
    ca-certificates \
    gnupg \
    lsb-release \
    openjdk-17-jre-headless \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js 24.x and npm
RUN curl -fsSL https://deb.nodesource.com/setup_24.x | bash - \
    && apt-get update && apt-get install -y nodejs \
    && npm install -g npm@latest

# Install Go 1.24
RUN wget -q https://go.dev/dl/go1.25.1.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz \
    && rm go1.25.1.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GO111MODULE=on

## Install Go tools (root). Use a temporary GOBIN so binary ends up in /usr/local/bin without keeping env var.
RUN GOBIN=/usr/local/bin go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Install Swagger CLI (optional, for validation)
RUN npm install -g @apidevtools/swagger-cli

WORKDIR /workspace

# Add user for VS Code
RUN useradd -m vscode && chown -R vscode /workspace
USER vscode

# Set per-user Go paths and install developer tools (gopls)
ENV GOPATH=/home/vscode/go
ENV GOBIN=/home/vscode/go/bin
ENV PATH=$GOBIN:$GOPATH/bin:$PATH
RUN go install golang.org/x/tools/gopls@latest

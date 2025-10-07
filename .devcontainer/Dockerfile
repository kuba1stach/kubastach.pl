# Use the base image from the devcontainer.json
FROM mcr.microsoft.com/devcontainers/base:ubuntu

# Install gcloud SDK
RUN apt-get update && \
    apt-get install -y curl apt-transport-https ca-certificates gnupg && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
    apt-get update && \
    apt-get install -y google-cloud-sdk

# Install additional features
RUN apt-get install -y git golang

# Install Node.js and npm from NodeSource
RUN curl -fsSL https://deb.nodesource.com/setup_current.x | bash - && \
    apt-get install -y nodejs

# Install React globally
RUN npm install -g typescript

# Clean up
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Set the default command to bash
CMD ["bash"]
# Use Ubuntu as the base image
FROM ubuntu:20.04

# Set non-interactive mode for apt-get
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    lsb-release \
    postgresql-client \
    postgresql \
    sudo \
    git \
    iptables \
    ca-certificates \
    apt-transport-https && \
    rm -rf /var/lib/apt/lists/*

# Ensure /etc/apt/keyrings directory exists (Fix for missing directory)
RUN mkdir -p /etc/apt/keyrings

# Install Docker CLI and Docker Compose Plugin
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo tee /etc/apt/keyrings/docker.asc > /dev/null && \
    echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && apt-get install -y \
    docker-ce-cli \
    docker-compose-plugin && \
    rm -rf /var/lib/apt/lists/*

# Install Go
RUN curl -LO https://golang.org/dl/go1.20.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz && \
    rm go1.20.4.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

# Ensure 'docker' group exists and add root user
RUN groupadd -g 998 docker || true
RUN usermod -aG docker root

# Set working directory
WORKDIR /workspace

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose port 8080 (for API) and 5432 (for PostgreSQL)
EXPOSE 8080 5432

# Start PostgreSQL and keep the container running
CMD ["bash", "-c", "service postgresql start && tail -f /dev/null"]

FROM golang:1.24-bookworm AS deps

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o main .

FROM debian:bookworm-slim

# Set DEBIAN_FRONTEND to noninteractive to avoid prompts during apt-get install
ENV DEBIAN_FRONTEND=noninteractive

# Set a working directory for temporary downloads
WORKDIR /tmp/setup

# Update package list and install base dependencies in one layer
# Also cleans up apt cache to reduce image size
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        git-all \
        wget \
        vim \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Download and install kubectl
# Fetches the latest stable version dynamically
RUN KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt) && \
    echo "Downloading kubectl version: ${KUBECTL_VERSION}" && \
    curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl && \
    kubectl version --client

# Download and install Helm
RUN echo "Downloading and installing Helm..." && \
    curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && \
    ./get_helm.sh && \
    rm get_helm.sh && \
    helm version

# Install K9s (v0.50.0 - as specified in the script)
RUN K9S_VERSION="v0.50.0" && \
    echo "Downloading K9s version: ${K9S_VERSION}" && \
    wget "https://github.com/derailed/k9s/releases/download/${K9S_VERSION}/k9s_Linux_amd64.tar.gz" && \
    # Extract only the k9s binary, license and readme might also be present
    tar -xzf k9s_Linux_amd64.tar.gz k9s && \
    chmod +x k9s && \
    mv k9s /usr/local/bin/ && \
    rm k9s_Linux_amd64.tar.gz && \
    k9s version

# Install Palette CLI (v4.6.1 - as specified in the script)
RUN PALETTE_VERSION="v4.6.1" && \
    echo "Downloading Palette CLI version: ${PALETTE_VERSION}" && \
    wget "https://software.spectrocloud.com/palette-cli/${PALETTE_VERSION}/linux/cli/palette" && \
    chmod +x palette && \
    mv palette /usr/local/bin/ && \
    palette version

# Install Palette Edge (v4.6.9 - as specified in the script)
RUN PALETTE_EDGE_VERSION="v4.6.9" && \
    echo "Downloading Palette Edge version: ${PALETTE_EDGE_VERSION}" && \
    wget "https://software.spectrocloud.com/stylus/${PALETTE_EDGE_VERSION}/cli/linux/palette-edge" && \
    chmod +x palette-edge && \
    mv palette-edge /usr/local/bin/ && \
    palette-edge --version

# Install Oras (v1.2.2 - as specified in the script)
# Ref: https://oras.land/docs/installation/
RUN ORAS_VERSION="1.2.2" && \
    echo "Downloading Oras version: v${ORAS_VERSION}" && \
    curl -LO "https://github.com/oras-project/oras/releases/download/v${ORAS_VERSION}/oras_${ORAS_VERSION}_linux_amd64.tar.gz" && \
    mkdir -p oras-install/ && \
    # Extract only the oras binary
    tar -zxf oras_${ORAS_VERSION}_linux_amd64.tar.gz -C oras-install/ oras && \
    mv oras-install/oras /usr/local/bin/ && \
    rm -rf oras_${ORAS_VERSION}_linux_amd64.tar.gz oras-install/ && \
    oras version

WORKDIR /app

COPY --from=builder /app/main .

COPY assets .

CMD ["/app/main", "webserver"]

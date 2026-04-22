FROM golang:1.24

WORKDIR /app

# Copy go module files and vendor directory
COPY go.mod go.sum ./
COPY vendor ./vendor

# Copy source code
COPY *.go ./

# Build the binary using vendored dependencies
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -mod=vendor -o setup .

# Copy static assets
COPY assets ./

# Set container defaults
ENV PORT=80
ENV SERVE_DIR=/app

# Default command starts the webserver
CMD ["/app/setup", "webserver"]

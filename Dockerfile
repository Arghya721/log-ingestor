# Use Ubuntu image for the builder stage
FROM ubuntu:latest AS builder

# Set the Go version
ENV GO_VERSION 1.20

# Install build dependencies and download Go
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    wget \
    librdkafka-dev \
 && wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz \
 && tar -xvf go${GO_VERSION}.linux-amd64.tar.gz \
 && mv go /usr/local

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Set the current working directory
WORKDIR /log-ingestor

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies, including confluent-kafka-go
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with CGO enabled
RUN go build -ldflags="-w -s" -o log-ingestor cmd/main.go

# Use Ubuntu latest for the final image
FROM ubuntu:latest

# Install runtime dependencies for Kafka (librdkafka)
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

# Set the current working directory
WORKDIR /log-ingestor

# Copy the binary from the builder stage
COPY --from=builder /log-ingestor/log-ingestor /log-ingestor

# Expose the necessary port
EXPOSE 1323

# Command to run the executable
CMD ["./log-ingestor"]

# Stage 1: Build the Go application
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY ./server ./server

# Build the Go application
WORKDIR /app/server
RUN go build -o main .

# Stage 2: Create the runtime environment
FROM debian:bookworm-slim

# Update apt cache separately to avoid re-downloading if package list doesn't change
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Install Chromium and X dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    xvfb \
    fonts-liberation && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Install miscellaneous libraries
RUN apt-get update && apt-get install -y --no-install-recommends \
    libasound2 \
    libatk1.0-0 \
    libatspi2.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xauth && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/server/main .

# Expose the application port
EXPOSE 8080

# Set environment variables for Chromium
ENV CHROMEDP_HEADLESS=true
ENV DISPLAY=:99

# Start xvfb and run the application
CMD ["sh", "-c", "xvfb-run -a ./main"]

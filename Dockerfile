# Builder stage
FROM golang:1.25.4-trixie AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install templ and build
RUN make dependencies
RUN make build

# Distroless final stage
FROM gcr.io/distroless/base-debian12:nonroot

# Copy the built application from builder stage
COPY --from=builder /app/bin/app /app

# Expose port (adjust if needed)
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app"]

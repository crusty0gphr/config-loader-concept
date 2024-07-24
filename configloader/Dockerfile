FROM golang:1.22-alpine AS config-loader-builder

WORKDIR /app

# Copy go.mod and go.sum files from the root directory
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY ../ ./

# Build the application
RUN go build -o /app/bin/configloader ./configloader/cmd

CMD ["/app/bin/configloader"]
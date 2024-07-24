FROM golang:1.22-alpine AS service-6b631a97-builder

WORKDIR /app

# Copy go.mod and go.sum files from the root directory
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY ../ ./

# Build the application
RUN go build -o /app/bin/service/6b631a97 ./servicepool/servcie6b631a97/cmd

CMD ["/app/bin/service/6b631a97"]

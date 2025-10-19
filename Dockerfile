# --- 1. Build Stage ---
# Use an official Go image as the base for building
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application source code
COPY . .

# Build the application. 
# CGO_ENABLED=0 creates a static binary (good for containers)
# -o /bin/server builds the output file to /bin/server
# ./cmd/server is the path to your main.go
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /bin/server ./cmd/server

# --- 2. Final Stage ---
# Use a minimal 'scratch' or 'alpine' image for the final container
# This makes your container small and secure
FROM alpine:latest

# Copy only the compiled binary from the 'builder' stage
COPY --from=builder /bin/server /bin/server

# (Optional) If your app needs SSL certificates
RUN apk --no-cache add ca-certificates

# Expose the port your Go app listens on
EXPOSE 8000 

# The command to run when the container starts
ENTRYPOINT ["/bin/server"]
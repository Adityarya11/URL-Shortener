# Go URL Shortener [🔗](https://url-shortener-liart-two.vercel.app/)

This is a high-performance URL shortener service built in Go, designed as an educational project to explore system design, API development, and database integration. It provides a simple API to create and resolve short URLs, with a pluggable repository layer supporting both in-memory and MongoDB persistent storage.

## Features

- **API Endpoints:** Provides `POST /shorten` to create links and `GET /{code}` to handle redirects.
- **Dual Storage Modes:** A toggleable persistence layer, configurable via an environment variable.
  - **In-Memory:** An ephemeral in-memory map for rapid development and testing.
  - **MongoDB:** A persistent MongoDB backend for production use.
- **Short Code Generation:** Supports both user-defined custom codes and random 6-character alphanumeric code generation.
- **Middleware:** Includes an in-memory, token-bucket rate limiter and configurable CORS middleware for API hardening.
- **Configuration:** All settings are managed via environment variables for easy setup and deployment.
- **Containerization:** Fully containerized with a multi-stage `Dockerfile` for efficient and repeatable production builds.

## Technology Stack

- **Backend:** Go (v1.22.5)
- **Router:** Standard Library `http.ServeMux`
- **Database:** MongoDB (using `go.mongodb.org/mongo-driver`)
- **Middleware:** `github.com/rs/cors` and a custom token-bucket rate limiter.
- **Deployment:** Docker (Future implementations!!)

## Architecture Overview

This project follows a clean, layered architecture to separate concerns and promote modularity.

- **Handlers (`internal/handlers`):** Responsible for parsing HTTP requests, validating input, and serializing JSON responses.
- **Services (`internal/services`):** Contains the core business logic, such as validating a URL, generating a short code, and checking for expiration.
- **Repository (`internal/repository`):** An interface-based layer (`Repository`) that abstracts all data storage operations. This allows the service layer to be agnostic of the database.
- **Main (`cmd/server`):** The application entry point responsible for initializing the database connection, injecting dependencies (like the `MongoRepo` or `MemoryRepo` into the `URLService`), setting up the router, and starting the HTTP server.

## Getting Started

### Prerequisites

- Go 1.22.5 or later
- Docker (for containerized deployment)
- A MongoDB Atlas account or a local MongoDB instance (if using `USE_MONGO=true`)

### Configuration

1.  Copy the example environment file:
    ```sh
    cp .env.example .env
    ```
2.  Edit the `.env` file with your settings:
    - `USE_MONGO`: Set to `true` to use MongoDB, or `false` to use the in-memory store.
    - `MONGO_URI`: Your MongoDB connection string (required if `USE_MONGO=true`).
    - `MONGO_DB`: Your MongoDB database name.
    - `RATE_LIMIT_REQUESTS`: Max requests per window (e.g., `100`).
    - `RATE_LIMIT_WINDOW`: Duration of the window (e.g., `1h`).

### Running Locally

1.  Install dependencies:
    ```sh
    go mod tidy
    ```
2.  Run the server (uses settings from `.env` file):
    ```sh
    go run ./cmd/server/main.go
    ```
    The server will start on port `8000` (or as specified by `$PORT`).

### Running with Docker (Just for flexing)

1.  Build the Docker image:
    ```sh
    docker build -t url-shortener .
    ```
2.  Run the container, passing in your environment file:
    ```sh
    docker run -p 8000:8000 --env-file .env url-shortener
    ```

## API Endpoints

### `POST /shorten`

Creates a new short URL.

**Request Body:**

```json
{
  "url": "[https://github.com/google/go-cmp](https://github.com/google/go-cmp)",
  "customCode": "go-cmp"
}
```

- `url` (string, required): The original URL to shorten.
- `customCode` (string, Optional): A desired custom short code. If not provided, a random 6-character code will be generated.

**Success Response (200 OK):**

```json
{
  "shortCode": "go-cmp",
  "shortUrl": "http://localhost:8000/go-cmp",
  "originalUrl": "https://github.com/google/go-cmp"
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "short code already exists"
}
```

### `GET /{code}`

Redirects a short code to its original URL.

- `code` (string, path parameter): The short code to resolve.

**Success Response (302 Found):**

- Redirects to the `OriginalURL` stored for the code.

- Increments the `ClickCount` for the link.

**Error Response (404 Not Found):**

- Returns if the code does not exist or has expired.

### `GET /health`

A health check endpoint to confirm the service is running.
**Success Response (200 OK):**

- `OK`

## Project Structure

```
url-shortener/
├── cmd/server/main.go     # Main application entry point and dependency injection.
├── internal/
│   ├── handlers/          # HTTP handlers (handler.go) and router setup (router.go).
│   ├── middleware/        # HTTP middleware (ratelimit.go).
│   ├── models/            # Core data structures (url.go).
│   ├── repository/        # Data storage layer:
│   │   ├── interface.go   # The core Repository interface.
│   │   ├── memory.go      # In-memory map implementation.
│   │   └── mongo.go       # MongoDB implementation.
│   └── services/          # Core business logic (url_service.go).
├── pkg/database/          # Database connection helpers (mongodb.go).
├── frontend/              # A simple static HTML/JS frontend for testing.
├── .env.example           # Example environment variables.
├── Dockerfile             # Multi-stage Docker build file.
├── go.mod                 # Go module dependencies.
└── Makefile               # Helper commands for development (build, test, run).
```

## Future Improvements

- Distributed Caching: Integrate a Redis cache (using the `pkg/database/redis.go` stub) to sit in front of the MongoDB repository. This will reduce database load for popular links and improve redirect latency.

- Advanced Analytics: Expand the `IncrementClicks` function to run in a goroutine, capturing and storing request metadata (User-Agent, Referrer, Geolocation) for a future analytics dashboard.

- Link Expiration: Implement TTL indexing in MongoDB to automatically purge expired links, or create a background job to handle cleanup based on the `ExpiresAt` field.

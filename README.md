# URL Shortener - MongoDB Edition

A production-ready URL shortener built from scratch in Go with MongoDB to learn system design principles and NoSQL patterns.

## Why MongoDB for URL Shortener?

- **Document-based storage**: URLs + metadata fit perfectly in JSON documents
- **Horizontal scaling**: Easy sharding across multiple servers
- **Flexible schema**: Add new fields without migrations
- **High write throughput**: Perfect for URL shortening workloads
- **JSON native**: No ORM complexity

## Features

- ✅ URL shortening with custom and generated codes
- ✅ MongoDB document storage with TTL
- ✅ Base62 encoding for short codes
- ✅ Analytics and click tracking
- ✅ In-memory caching (Phase 1) → Redis (Phase 3)
- ✅ Rate limiting and CORS middleware
- ✅ Comprehensive testing
- ✅ Production-ready logging

## Quick Start

### Prerequisites

- Go 1.21+
- MongoDB 7.0+ (installed locally)

### MongoDB Installation

**macOS (Homebrew)**:

```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb/brew/mongodb-community
```

**Ubuntu/Debian**:

```bash
# Import MongoDB public GPG Key
curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor

# Create list file
echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list

# Install MongoDB
sudo apt-get update
sudo apt-get install -y mongodb-org

# Start MongoDB
sudo systemctl start mongod
sudo systemctl enable mongod
```

**Windows**:
Download MongoDB Community Server from [official website](https://www.mongodb.com/try/download/community)

### Development Setup

1. **Clone and setup**:

   ```bash
   git clone <your-repo>
   cd url-shortener
   make setup
   ```

2. **Start MongoDB**:

   ```bash
   make mongo-start
   ```

3. **Initialize Go module**:

   ```bash
   go mod init url-shortener
   go mod tidy
   ```

4. **Start the server**:
   ```bash
   make run
   ```

The API will be available at `http://localhost:8000`

### API Usage

**Shorten a URL**:

```bash
curl -X POST http://localhost:8000/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "customCode": "example"}'
```

**Response**:

```json
{
  "shortCode": "example",
  "shortUrl": "http://localhost:8000/example",
  "originalUrl": "https://example.com",
  "expiresAt": "2026-08-28T00:00:00Z"
}
```

**Access shortened URL**:

```bash
curl -L http://localhost:8000/example
# Redirects to https://example.com
```

**Get analytics**:

```bash
curl http://localhost:8000/analytics/example
```

## MongoDB Document Structure

```json
{
  "_id": ObjectId("..."),
  "shortCode": "abc123",
  "originalUrl": "https://example.com",
  "customCode": true,
  "clickCount": 42,
  "createdAt": ISODate("2025-08-28T00:00:00Z"),
  "expiresAt": ISODate("2026-08-28T00:00:00Z"),
  "userId": "user123",
  "metadata": {
    "title": "Example Website",
    "description": "An example website"
  },
  "analytics": {
    "countries": {"US": 20, "IN": 15, "UK": 7},
    "browsers": {"Chrome": 30, "Safari": 12},
    "referrers": ["google.com", "facebook.com"]
  }
}
```

## Project Phases

### Phase 1: In-Memory Storage ✅

- Basic HTTP server with `net/http`
- In-memory map for URLs
- Base62 short code generation
- JSON API endpoints

### Phase 2: MongoDB Integration 🔄

- MongoDB connection and operations
- Document-based URL storage
- TTL (Time To Live) for expiration
- Basic indexing strategy

### Phase 3: Redis Caching 📋

- Redis for frequently accessed URLs
- Cache-aside pattern
- Performance optimization

### Phase 4: Advanced Features 📋

- Rate limiting middleware
- Bulk operations
- Advanced analytics
- Custom domains

### Phase 5: Production Scaling 📋

- MongoDB replica sets
- Sharding strategies
- Performance monitoring
- Load testing

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│   Go API    │───▶│  Service    │
└─────────────┘    │   Server    │    │   Layer     │
                   └─────────────┘    └─────────────┘
                           │                 │
                           ▼                 ▼
                   ┌─────────────┐    ┌─────────────┐
                   │    Redis    │    │  MongoDB    │
                   │   (Cache)   │    │ (Database)  │
                   └─────────────┘    └─────────────┘
```

## Development Commands

```bash
# Setup development environment
make setup

# Start/stop MongoDB
make mongo-start
make mongo-stop
make mongo-status

# Development
make run          # Run the server
make dev          # Run with hot reload (requires air)
make test         # Run tests
make coverage     # Test coverage
make lint         # Lint code

# MongoDB operations
make mongo-shell  # Connect to MongoDB shell

# Build
make build        # Development build
make build-prod   # Production build
```

## MongoDB Advantages for This Project

### 1. **Natural Document Structure**

URLs and their metadata fit perfectly in JSON documents - no complex relational mapping needed.

### 2. **Horizontal Scaling**

```javascript
// Easy sharding by shortCode
sh.shardCollection("urlshortener.urls", { shortCode: 1 });
```

### 3. **Flexible Analytics**

```javascript
// Easy to add new analytics fields
db.urls.updateOne({ shortCode: "abc123" }, { $inc: { "analytics.mobile": 1 } });
```

### 4. **TTL Collections**

```javascript
// Automatic URL expiration
db.urls.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 });
```

### 5. **Aggregation Pipelines**

```javascript
// Complex analytics queries
db.urls.aggregate([
  { $match: { createdAt: { $gte: ISODate("2025-01-01") } } },
  { $group: { _id: "$userId", totalClicks: { $sum: "$clickCount" } } },
]);
```

## Performance Targets

- **Latency**: <50ms for shortening, <10ms for redirection
- **Throughput**: 1000+ requests/second
- **Storage**: 10M+ URLs with efficient indexing
- **Cache Hit Rate**: 90%+ for popular URLs

## Contributing

This is an educational project for learning system design. Feel free to:

1. Fork and experiment
2. Add new features
3. Optimize performance
4. Add more comprehensive tests

## Learning Resources

- [MongoDB Go Driver Documentation](https://pkg.go.dev/go.mongodb.org/mongo-driver)
- [MongoDB Best Practices](https://www.mongodb.com/developer/products/mongodb/mongodb-schema-design-best-practices/)
- [System Design Primer](https://github.com/donnemartin/system-design-primer)

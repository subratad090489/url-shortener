
# URL Shortener Service - Modular Architecture

A modular URL shortening service built with Go, following clean architecture principles.

## Project Structure

```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── handlers/
│   │   └── handlers.go            # HTTP handlers
│   │   └── handlers_test.go       # Handler tests
│   ├── models/
│   │   └── models.go              # Data models
│   └── shortener/
│       └── shortener.go           # Business logic
│       └── shortener_test.go         # ← Test file location
├── deployments/
│   └── docker/
│       ├── Dockerfile             # Docker build file
│       └── docker-compose.yml     # Docker Compose configuration
├── go.mod                         # Go module file
├── Makefile                       # Build automation
└── README.md                      # Documentation
```

# Architecture

### Layers

1. **cmd/server** - Application entry point and initialization
2. **internal/handlers** - HTTP request handling and routing
3. **internal/shortener** - Core business logic for URL shortening
4. **internal/models** - Data transfer objects and models
5. **internal/config** - Configuration management
6. **deployments/docker** - Docker and deployment configurations

## API Endpoints

### 1. Shorten URL
**POST** `/shorten`
```json
{
  "long_url": "https://www.example.com/very/long/url"
}
```

### 2. Redirect
**GET** `/{shortCode}` - Returns HTTP 301 redirect

### 3. Health Check
**GET** `/health` - Returns service health status

### 4. Statistics
**GET** `/stats` - Returns service statistics

## Running the Service

### Local Development
```bash
# Or build and run
make build
./bin/url-shortener
```
### Using Docker
```bash
# Build and run with Docker
make docker-build
make docker-run
```

### Using Docker Compose (Recommended)
```bash
# Start the service
make docker-compose-up

# Stop the service
make docker-compose-down
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `HOST` - Server host (default: localhost)


## Usage Examples

```bash
# Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://www.google.com"}'

# Get statistics
curl http://localhost:8080/stats

# Health check
curl http://localhost:8080/health

# Use short URL (will redirect)
curl -L http://localhost:8080/{shortCode}
```

## Development

### Running Tests
```bash
make test
```

### Building
```bash
make build
```

### Cleaning Up
```bash
make clean
```

## To deploy with a single command:
```bash
make docker-compose-up
```

### Limitations/Improvements
1. **Current code doesn't have support for HTTPs request**
2. **It supports only incoming request content-type application/json**
3. **Used in-mermory datastructure for storing data**

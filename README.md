<<<<<<< HEAD
# beto
Go project enviroment setting up
=======
# Beto Application

A modern, production-ready Go web application with comprehensive development tooling and best practices.

## Features

- 🚀 **High Performance**: Built with Go for optimal performance and concurrency
- 🔧 **Developer Experience**: Live reload, comprehensive testing, and debugging tools
- 🐳 **Docker Support**: Multi-stage builds and Docker Compose for development
- 📊 **Monitoring**: Health checks, metrics, and logging infrastructure
- 🔒 **Security**: Built-in security best practices and middleware
- 🧪 **Testing**: Unit tests, benchmarks, and test coverage
- 📦 **CI/CD Ready**: Automated builds, testing, and deployment pipelines
- 🔍 **Code Quality**: Linting, formatting, and static analysis tools

## Quick Start

### Prerequisites

- Go 1.25 or later
- Docker and Docker Compose (optional)
- Make (optional, but recommended)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/darkcloud/beto.git
   cd beto
   ```

2. **Setup development environment**
   ```bash
   make setup
   ```

3. **Copy environment configuration**
   ```bash
   make env
   ```

4. **Run the application**
   ```bash
   make run
   ```

The application will be available at `http://localhost:8080`

### Development with Live Reload

For the best development experience with automatic reloading:

```bash
make dev
```

## Project Structure

```
.
├── api/                    # API definitions and OpenAPI specs
├── cmd/                    # Application entry points
├── configs/                # Configuration files
├── docs/                   # Documentation
├── internal/               # Private application code
├── pkg/                    # Public library code
│   ├── config/            # Configuration management
│   └── logger/            # Structured logging
├── scripts/               # Build and deployment scripts
├── test/                  # Integration tests
├── .air.toml             # Live reload configuration
├── .env.example          # Environment variables template
├── .golangci.yml         # Linter configuration
├── docker-compose.yml    # Docker services configuration
├── Dockerfile            # Production Docker image
├── Dockerfile.dev        # Development Docker image
├── go.mod                # Go module dependencies
├── main.go               # Application entry point
├── main_test.go          # Application tests
└── Makefile             # Build automation
```

## Development

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make dev               # Run with live reload
make test              # Run all tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make fmt               # Format code
make check             # Run all code quality checks
```

### Environment Configuration

Copy `.env.example` to `.env` and configure the following variables:

```bash
# Application
PORT=8080
APP_NAME=Beto Application
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=beto_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Testing

Run the complete test suite:

```bash
make test              # Unit tests
make test-coverage     # Coverage report
make test-race         # Race condition detection
make bench            # Benchmarks
```

### Code Quality

Maintain code quality with:

```bash
make check            # All quality checks
make lint             # Linting
make fmt              # Code formatting
make vet              # Go vet analysis
make staticcheck      # Static analysis
```

## API Endpoints

### Health and Status

- `GET /health` - Health check endpoint
- `GET /version` - Application version information
- `GET /` - Root endpoint with welcome message

### API v1

- `GET /api/v1/status` - API status and uptime information

### Example Responses

**Health Check:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Version:**
```json
{
  "name": "Beto Application",
  "version": "1.0.0"
}
```

## Docker Support

### Development with Docker

Run the entire stack with Docker Compose:

```bash
# Basic stack (app, postgres, redis)
docker-compose up -d

# Development mode with live reload
docker-compose --profile dev up -d

# With monitoring (Prometheus, Grafana)
docker-compose --profile monitoring up -d

# With logging (ELK stack)
docker-compose --profile logging up -d
```

### Production Deployment

Build and run the production Docker image:

```bash
make docker-build
make docker-run
```

## Configuration

### Application Configuration

The application uses a hierarchical configuration system:

1. **Default values** - Built into the application
2. **Environment variables** - Override defaults
3. **Configuration files** - For complex configurations

### Logging

Structured logging with configurable levels and formats:

- **Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Formats**: JSON, Text
- **Context**: Request ID, User ID, and custom fields

```go
import "github.com/darkcloud/beto/pkg/logger"

log := logger.WithFields(map[string]interface{}{
    "user_id": "12345",
    "action": "login",
})
log.Info("User logged in successfully")
```

## Monitoring and Observability

### Health Checks

- Application health: `GET /health`
- Kubernetes-ready health checks
- Docker health check support

### Metrics (Optional)

When monitoring profile is enabled:

- **Prometheus**: Metrics collection at `:9090`
- **Grafana**: Dashboards at `:3000` (admin/admin)

### Logging (Optional)

When logging profile is enabled:

- **Elasticsearch**: Search and analytics at `:9200`
- **Kibana**: Log visualization at `:5601`
- **Logstash**: Log processing at `:5000`

## Security

### Built-in Security Features

- **CORS**: Configurable cross-origin resource sharing
- **Request Timeout**: Prevents slow attacks
- **Graceful Shutdown**: Proper connection handling
- **Environment Variables**: Secure configuration management

### Best Practices

- Never commit secrets to version control
- Use environment variables for sensitive data
- Regularly update dependencies
- Run security scans with `make security`

## Performance

### Benchmarking

Run performance benchmarks:

```bash
make bench
```

### Profiling

CPU and memory profiling:

```bash
make profile-cpu      # CPU profiling
make profile-mem      # Memory profiling
```

## Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Run tests**: `make test`
4. **Run quality checks**: `make check`
5. **Commit changes**: `git commit -m 'Add amazing feature'`
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Open a Pull Request**

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` and `goimports` for formatting
- Write tests for new functionality
- Update documentation as needed

## Deployment

### Environment Setup

1. **Staging Environment**
   ```bash
   APP_ENV=staging make build
   ```

2. **Production Environment**
   ```bash
   APP_ENV=production make build
   ```

### Build Targets

```bash
make build             # Linux binary
make build-windows     # Windows binary
make build-mac         # macOS binary
make build-all         # All platforms
```

### Release Process

```bash
make release          # Complete release build
make tag              # Create git tag
```

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Check what's using port 8080
lsof -i :8080

# Use a different port
PORT=8081 make run
```

**Database connection failed:**
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View database logs
docker-compose logs postgres
```

**Redis connection failed:**
```bash
# Check Redis status
docker-compose ps redis

# Test Redis connection
docker-compose exec redis redis-cli ping
```

### Debug Mode

Enable debug logging:

```bash
LOG_LEVEL=debug make run
```

### Getting Help

- Check the [documentation](./docs/)
- Review [example configurations](./configs/)
- Open an [issue](https://github.com/darkcloud/beto/issues)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- HTTP routing by [Gorilla Mux](https://github.com/gorilla/mux)
- Live reload with [Air](https://github.com/air-verse/air)
- Linting by [golangci-lint](https://github.com/golangci/golangci-lint)
- Testing with [Testify](https://github.com/stretchr/testify)

---

**Happy coding! 🚀**
>>>>>>> b25d782 (initial commit)

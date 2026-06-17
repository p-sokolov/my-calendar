# Calendar API
![CI](https://github.com/p-sokolov/my-calendar/actions/workflows/ci.yml/badge.svg)

REST API service for calendar event management written in Go.

## Features

- CRUD operations for calendar events
- OpenAPI/Swagger documentation
- Prometheus metrics collection
- Grafana monitoring dashboards
- ClickHouse request analytics
- Structured logging with Zap
- Docker & Docker Compose support
- Postman collection for API testing
- GitHub Actions CI pipeline
- Static code analysis with golangci-lint

## Tech Stack

- Go
- Gorilla Mux
- ClickHouse
- Prometheus
- Grafana
- Docker
- Swagger/OpenAPI
- Postman
- GitHub Actions

## Architecture

```text
Client
   │
   ▼
Calendar API
   │
   ├── Prometheus ──► Grafana
   │
   └── ClickHouse ──► Grafana
```

## Run

```bash
docker compose up --build
```

### Endpoints

API:

```text
http://localhost:8080
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

Grafana:

```text
http://localhost:3000
```

Prometheus:

```text
http://localhost:9090
```

## Monitoring

Prometheus metrics:

- HTTP requests count
- Request duration histogram
- Endpoint-level metrics

Grafana dashboards:

- Requests per second
- Response latency
- Endpoint analytics
- Status code distribution

## Testing

```bash
go test ./...
```

```bash
golangci-lint run
```

Postman collection:

```text
postman/calendar-api.json
```

## CI

GitHub Actions pipeline runs:

- unit tests
- golangci-lint checks

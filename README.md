
![CI](https://github.com/p-sokolov/my-calendar/actions/workflows/ci.yml/badge.svg)

# Calendar API

A REST API service for managing calendar events.

The service provides CRUD operations for events and demonstrates a production-style backend setup with monitoring, analytics, API documentation, testing, containerization and CI/CD automation.

---

# Task Scope and Expectations

The goal of the project is to implement a simple calendar service while applying commonly used backend engineering practices and tooling.

The project focuses on:

- REST API development in Go
- API documentation with Swagger
- Monitoring with Prometheus and Grafana
- Request analytics with ClickHouse
- Containerized deployment with Docker Compose
- Automated testing and linting
- CI pipeline using GitHub Actions

---

# Functional Requirements

The service supports:

- Create event
- Read events
- Update event
- Delete event

Available endpoints:

| Method | Endpoint |
|----------|----------|
| GET | /events |
| POST | /events |
| PUT | /events/{id} |
| DELETE | /events/{id} |

Additional endpoints:

| Endpoint | Purpose |
|----------|----------|
| /metrics | Prometheus metrics |
| /swagger/index.html | Swagger UI |

---

# Technical Requirements

### Backend

- Go
- Gorilla Mux
- Zap Logger

### Documentation

- Swagger / OpenAPI

### Monitoring

- Prometheus
- Grafana

### Analytics

- ClickHouse

### Infrastructure

- Docker
- Docker Compose

### Quality Assurance

- Unit Tests
- golangci-lint
- GitHub Actions CI

### API Testing

- Postman Collection

---

# Architecture

```text
                     +----------------+
                     |     Client     |
                     +--------+-------+
                              |
                              v
                    +------------------+
                    |   Calendar API   |
                    |       Go         |
                    +---+----------+---+
                        |          |
            Metrics     |          | Request Logs
                        |          |
                        v          v
               +---------------+  +---------------+
               |  Prometheus   |  |  ClickHouse   |
               +-------+-------+  +-------+-------+
                       |                  |
                       +--------+---------+
                                |
                                v
                         +-------------+
                         |   Grafana   |
                         +-------------+
```

---

# How to Run

Start all services:

```bash
docker compose up --build
```

Available services:

| Service | URL |
|----------|----------|
| API | http://localhost:8080 |
| Swagger | http://localhost:8080/swagger/index.html |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| ClickHouse HTTP | http://localhost:8123 |

---

# Testing

Run unit tests:

```bash
go test ./...
```

Run linter:

```bash
golangci-lint run
```

Postman collection:

```text
postman/calendar-api.json
```

---

# Solution Notes

This project intentionally keeps business logic simple while focusing on backend engineering practices often used in production systems.

Implemented additions include:

- Prometheus metrics collection
- Grafana dashboards
- ClickHouse request analytics
- Structured logging
- Dockerized environment
- Swagger documentation
- Postman collection
- GitHub Actions CI pipeline
- Static code analysis with golangci-lint

The project serves as a demonstration of building, monitoring, testing and maintaining a Go REST API service.


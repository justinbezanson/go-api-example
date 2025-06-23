# Roller Coaster API in Go

## Getting Started

### Prerequisites

- Go 1.24.3

## Running the Application

```bash
ADMIN_PASSWORD=admin go run server.go
```

## CURL commands

- GET /coasters

```bash
curl -X GET http://localhost:8080/coasters
```

- GET /coasters/:id

```bash
curl -X GET http://localhost:8080/coasters/1
```

- POST /coasters

```bash
curl -X POST http://localhost:8080/coasters \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Roller Coaster 1",
    "manufacturer": "Coaster Inc.",
    "id": "1",
    "in_park": "Park A",
    "height": 50
  }'
 ```

 - POST /coasters/random

 ```bash
 curl -X POST http://localhost:8080/coasters/random -L
 ```

 - GET /admin
 ```bash
 curl http://localhost:8080/admin -X GET -u admin:admin
 ```
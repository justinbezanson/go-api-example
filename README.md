# go-api-example

## Getting Started

### Prerequisites

- Go 1.24.3

## Running the Application

```bash
go run server.go
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
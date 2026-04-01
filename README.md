# Sample ATG App

Go + MongoDB sample service with CRUD flows for users, addresses, products, orders, and payments.

## Run locally

```bash
docker compose up -d --build
```

The API listens on `http://localhost:8080`.

## Run the smoke flow

```bash
docker compose --profile test run --rm tester
```

## API reference

See [API.md](./API.md) for the schema and curl examples for every endpoint.

## CI

GitHub Actions runs on pull requests targeting `main` and executes:

1. `go test ./...`
2. the Docker Compose smoke flow through the tester container

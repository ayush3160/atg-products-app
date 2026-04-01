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

## Run on Kubernetes

Build the API image locally first:

```bash
docker build -t sample-atg-app-api:latest .
```

If you are using `kind` or `minikube`, load that image into the cluster before applying the manifests.

```bash
kind load docker-image sample-atg-app-api:latest
# or
minikube image load sample-atg-app-api:latest
```

Apply the Kubernetes manifests:

```bash
kubectl apply -k k8s
```

Access the API with port-forwarding:

```bash
kubectl -n sample-atg port-forward svc/api 8080:80
```

The API then listens on `http://localhost:8080`.

## API reference

See [API.md](./API.md) for the schema and curl examples for every endpoint.

## CI

GitHub Actions runs on pull requests targeting `main` and executes:

1. `go test ./...`
2. the Docker Compose smoke flow through the tester container

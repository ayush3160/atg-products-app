FROM golang:1.25.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/sample-atg-app ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/sample-atg-app /usr/local/bin/sample-atg-app

USER app
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/sample-atg-app"]

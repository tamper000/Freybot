# Stage 1
FROM golang:1.25.1-bookworm as builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libssl-dev \
    ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/main.go

RUN mkdir -p /src/config /src/database

# Stage 2
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

RUN adduser --disabled-password --gecos '' appuser

COPY --from=builder /src/main /app/main
COPY --from=builder /src/config /app/config
COPY --from=builder /src/database /app/database

WORKDIR /app

RUN chown -R appuser:appuser /app; \
    mkdir -p /src/config /src/database

USER appuser

# webhook
EXPOSE 8888
# metrics
EXPOSE 8080

ENTRYPOINT ["./main"]

# Stage 1
FROM golang:1.24.5-bookworm as builder

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

RUN mkdir -p /src/config; \
    mkdir -p /src/database; \
    touch /src/database/bot.db

# Stage 2
FROM gcr.io/distroless/cc-debian12

COPY --chown=65532:65532 --from=builder /src/main /app/main
COPY --chown=65532:65532 --from=builder /src/config /app/config
COPY --chown=65532:65532 --from=builder /src/database /app/database

WORKDIR /app
USER 65532

EXPOSE 8888
ENTRYPOINT ["./main"]

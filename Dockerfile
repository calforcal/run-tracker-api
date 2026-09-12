FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/run-tracker-api ./cmd

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=build /out/run-tracker-api ./run-tracker-api
COPY --from=build /src/internal/adapters/postgres/migrations ./internal/adapters/postgres/migrations

USER appuser

ENV PORT=8000
EXPOSE 8000

ENTRYPOINT ["./run-tracker-api"]

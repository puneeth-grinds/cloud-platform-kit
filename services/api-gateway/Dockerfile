# --- Stage 1: Build Stage ---

FROM golang:1.27.2-alpine3.24 AS build

WORKDIR /app

COPY go.mod .

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/api-gateway ./cmd


# --- Stage 2: Production Stage ---

FROM alpine:3.24 AS runner

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN addgroup -g 10001 app && adduser -D -u 10001 -G app appuser

COPY --from=build --chown=appuser:app app/api-gateway .

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/api-gateway"]
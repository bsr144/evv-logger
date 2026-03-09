FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETARCH

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /app/bin/evv-logger ./cmd/http

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/bin/evv-logger /app/evv-logger

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/evv-logger"]

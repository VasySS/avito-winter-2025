##################################################
# Build stage
##################################################
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main ./cmd/server/main.go

##################################################
# Runtime stage
##################################################
FROM alpine:3.21 AS runtime

WORKDIR /app
COPY --from=builder build/main .

RUN apk add --no-cache ca-certificates curl && update-ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

ARG SERVER_PORT
ENV SERVER_PORT $SERVER_PORT
HEALTHCHECK --start-period=10s CMD curl -f "http://localhost:${SERVER_PORT}/health" || exit 1

EXPOSE $SERVER_PORT
ENTRYPOINT ["./main"]

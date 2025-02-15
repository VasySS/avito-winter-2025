##################################################
# Build stage
##################################################
FROM golang:1.24-alpine:3.21 AS builder

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

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD [ "curl", "-f", "http://localhost:4174/health" ]

EXPOSE 8080
ENTRYPOINT ["./main"]

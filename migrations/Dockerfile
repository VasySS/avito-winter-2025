# https://github.com/kukymbr/goose-docker/blob/main/Dockerfile
FROM alpine:3.21

ENV GOOSE_VERSION_TAG="v3.24.1"

RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

ARG TARGETARCH
RUN case "$TARGETARCH" in \
      "amd64"|"") GOOSE_ARCH="x86_64";; \
      "arm64") GOOSE_ARCH="arm64";; \
      *) echo "Unsupported architecture: $TARGETARCH; switching to amd64."; GOOSE_ARCH="x86_64";; \
    esac && \
    echo "Requesting goose binary $GOOSE_VERSION_TAG for $GOOSE_ARCH (TARGETARCH=$TARGETARCH)..." && \
    wget -O /bin/goose "https://github.com/pressly/goose/releases/download/$GOOSE_VERSION_TAG/goose_linux_$GOOSE_ARCH"

WORKDIR /migrations
COPY . .

RUN chmod +x /bin/goose
ENTRYPOINT [ "/bin/goose"]
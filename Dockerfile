ARG ALPINE_VERSION=3.11
ARG GO_VERSION=1.14

FROM alpine:${ALPINE_VERSION} AS alpine
RUN apk --update add tzdata

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
RUN apk --update add git libcap
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY cmd/pingodown/main.go ./
COPY internal internal
RUN go build -ldflags="-s -w" -o app main.go
RUN setcap cap_net_raw=+ep /tmp/gobuild/app

FROM scratch
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
LABEL \
    org.opencontainers.image.authors="quentin.mcgaw@gmail.com" \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.url="https://github.com/qdm12/pingodown" \
    org.opencontainers.image.documentation="https://github.com/qdm12/pingodown/blob/master/README.md" \
    org.opencontainers.image.source="https://github.com/qdm12/pingodown" \
    org.opencontainers.image.title="pingodown" \
    org.opencontainers.image.description="Introduce more ping on a UDP port, for gaming purposes"
COPY --from=alpine --chown=1000 /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=America/Montreal \
    LISTEN_ADDRESS=:8000 \
    SERVER_ADDRESS=
ENTRYPOINT ["/app"]
USER 1000
COPY --from=builder --chown=1000 /tmp/gobuild/app /app

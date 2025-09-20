FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT_SHA
ARG BUILD_TIME

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build \
        -ldflags="-X github.com/alexraskin/LhBotGo/internal/ver.buildVersion=$VERSION -X github.com/alexraskin/LhBotGo/internal/ver.buildCommit=$COMMIT_SHA -X github.com/alexraskin/LhBotGo/internal/ver.buildTime=$BUILD_TIME" \
        -o lhbotgo github.com/alexraskin/LhBotGo

FROM alpine

COPY --from=build /build/lhbotgo /bin/lhbotgo

ENTRYPOINT ["/bin/lhbotgo"]

CMD ["-config", "/var/lib/config.toml"]
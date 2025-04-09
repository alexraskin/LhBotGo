FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o lhbotgo github.com/alexraskin/LhBotGo

FROM alpine

COPY --from=build /build/lhbotgo /bin/lhbotgo

ENTRYPOINT ["/bin/lhbotgo"]

CMD ["-config", "/var/lib/config.toml"]
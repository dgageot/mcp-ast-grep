ARG GO_VERSION=1.25.0

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS base
WORKDIR /app

FROM base AS mcp-server
RUN apk add --no-cache ast-grep
ARG TARGETOS
ARG TARGETARCH
RUN --mount=target=.\
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags "-s -w" -o /mcp-server .
ENTRYPOINT [ "/mcp-server" ]
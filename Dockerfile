ARG GO_VERSION=1.25.0

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS base
WORKDIR /app

FROM base AS mcp-server
RUN apk add --no-cache ast-grep
ARG TARGETOS
ARG TARGETARCH
LABEL io.docker.server.metadata="{"name": "ast-grep", "volumes": [\"{{ast-grep.path|volume-target}}:/src\"], "config": [{"name": "ast-grep", "type": "object", "properties": {"path": {"type": "string"}}, "required": ["path"]}]}"
RUN --mount=target=.\
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags "-s -w" -o /mcp-server .
WORKDIR /src
ENTRYPOINT [ "/mcp-server" ]
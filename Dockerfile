ARG GO_VERSION=1.25.1

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS build-mcp-server
WORKDIR /app
ARG TARGETOS
ARG TARGETARCH
LABEL io.docker.server.metadata="{"name": "ast-grep", "volumes": [\"{{ast-grep.path|volume-target}}:/src\"], "config": [{"name": "ast-grep", "type": "object", "properties": {"path": {"type": "string"}}, "required": ["path"]}]}"
COPY . ./
ADD https://raw.githubusercontent.com/ast-grep/ast-grep-mcp/b69eb5391bd93d46ef3dec07de814c3c39675c8f/ast-grep.mdc ./instructions.md
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags "-s -w" -o /mcp-server .

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1
RUN apk add --no-cache ast-grep
WORKDIR /src
COPY --from=build-mcp-server /mcp-server /
ENTRYPOINT ["/mcp-server"]
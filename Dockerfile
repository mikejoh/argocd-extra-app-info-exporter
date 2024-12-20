# syntax = docker/dockerfile:1.3
ARG GO_VERSION=1.23.4

FROM golang:$GO_VERSION AS base
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

FROM base AS build
ARG VERSION
ARG GIT_SHA
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build \
    -o /build/argocd-extra-app-info-exporter \
    -ldflags "-s \
    -w \
    -X=github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo.Version=$VERSION \
    -X=github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo.Name=argocd-extra-app-info-exporter \
    -X=github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo.GitSHA=$GIT_SHA" \
    ./cmd/argocd-extra-app-info-exporter/*.go

FROM gcr.io/distroless/static:nonroot
COPY --from=build /build/argocd-extra-app-info-exporter /bin/argocd-extra-app-info-exporter
USER 65532:65532
ENTRYPOINT ["/bin/argocd-extra-app-info-exporter"]
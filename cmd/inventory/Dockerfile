# syntax = docker/dockerfile:1-experimental

# hadolint ignore=DL3007
FROM jacodelg/base:latest AS base

WORKDIR /go/src/github.com/jacob-delgado/playground
COPY . /go/src/github.com/jacob-delgado/playground

RUN /usr/bin/buf generate && \
    go mod download

WORKDIR /go/src/github.com/jacob-delgado/playground
COPY . /go/src/github.com/jacob-delgado/playground

# hadolint ignore=DL3007
FROM jacodelg/base:latest AS buf

RUN /usr/bin/buf generate && \
    go mod download

RUN /usr/bin/buf breaking --against ".git#branch=main" && \
    /usr/bin/buf lint

FROM base AS unit-test
RUN --mount=type=cache,target=/root/.cache/go-build go test -v .

FROM golangci/golangci-lint:${GOLANGCI_LINT} AS lint-base

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/.cache/golangci-lint \
  golangci-lint run --timeout 10m0s ./...

FROM base AS build
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -o /out/inventory ./cmd/inventory

FROM scratch as inventory
COPY --from=build /out/inventory /
CMD ["/inventory"]

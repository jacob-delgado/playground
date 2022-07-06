# syntax = docker/dockerfile:1-experimental

ARG GOLANGCI_LINT=v1.46.2
FROM golang:1.18-bullseye as base

ENV BUF_BUILD=v1.6.0
ENV PROTOC_GEN_GO=v1.28.0
ENV PROTOC_GEN_GO_GRPC=v1.2.0
ENV PROTOC_GEN_VALIDATE=v0.6.7
ENV GOLANGCI_LINT=v1.46.2

RUN wget -nv -O /usr/bin/buf "https://github.com/bufbuild/buf/releases/download/${BUF_BUILD}/buf-Linux-$(uname -m)" && \
    chmod 555 /usr/bin/buf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC}
RUN go install github.com/envoyproxy/protoc-gen-validate@${PROTOC_GEN_VALIDATE}
RUN export PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /go/src/github.com/jacob-delgado/playground
ADD . /go/src/github.com/jacob-delgado/playground

RUN /usr/bin/buf generate
RUN go mod download

FROM base AS buf

RUN /usr/bin/buf breaking --against ".git#branch=main"
RUN /usr/bin/buf lint

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

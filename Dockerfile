# Start by building the application.
FROM golang:1.18-bullseye as build

ENV BUF_BUILD=v1.6.0
ENV PROTOC_GEN_GO=v1.28.0
ENV PROTOC_GEN_GO_GRPC=v1.2.0
ENV PROTOC_GEN_VALIDATE_VERSION=v0.6.7

RUN wget -nv -O /usr/bin/buf "https://github.com/bufbuild/buf/releases/download/${BUF_BUILD}/buf-Linux-$(uname -m)" && \
    chmod 555 /usr/bin/buf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC}
RUN go install github.com/envoyproxy/protoc-gen-validate@${PROTOC_GEN_VALIDATE_VERSION}
RUN export PATH="$PATH:$(go env GOPATH)/bin"

WORKDIR /go/src/github.com/jacob-delgado/playground
ADD . /go/src/github.com/jacob-delgado/playground

RUN /usr/bin/buf generate
RUN GOOS=linux go build ./cmd/playground

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian11
COPY --from=build /go/src/github.com/jacob-delgado/playground/playground /
CMD ["/playground"]

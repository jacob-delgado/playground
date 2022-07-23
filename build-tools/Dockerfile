ARG GOLANGCI_LINT=v1.46.2
FROM golang:1.18-bullseye as base

ENV BUF_BUILD=v1.6.0
ENV GOFUMPT=v0.3.1
ENV GOLANGCI_LINT=v1.46.2
ENV HADOLINT=v2.10.0
ENV KIND=v0.14.0
ENV KUBECTL=v1.24.2
ENV PROTOC_GEN_GO=v1.28.0
ENV PROTOC_GEN_GO_GRPC=v1.2.0
ENV PROTOC_GEN_VALIDATE=v0.6.7
ENV SHELLCHECK=v0.8.0
ENV TRIVY=0.28.1

# install docker
#hadolint ignore=DL3009
RUN apt-get update && \
    apt-get remove docker docker.io containerd runc && \
    apt-get autoremove
#hadolint ignore=DL3008
RUN apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    xz-utils
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
#hadolint ignore=DL3008,DL3009
RUN mkdir -p /etc/apt/keyrings && \
    wget -q -O - https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
      $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && \
    apt-get -y --no-install-recommends install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# buf
RUN wget -nv -O /usr/bin/buf "https://github.com/bufbuild/buf/releases/download/${BUF_BUILD}/buf-Linux-$(uname -m)" && \
    chmod 555 /usr/bin/buf

# hadolint
ADD https://github.com/hadolint/hadolint/releases/download/${HADOLINT}/hadolint-Linux-x86_64 /usr/bin/hadolint
RUN chmod 555 /usr/bin/hadolint

# kubectl
ADD https://storage.googleapis.com/kubernetes-release/release/${KUBECTL}/bin/linux/amd64/kubectl /usr/bin/kubectl
RUN chmod 555 /usr/bin/kubectl

# shellcheck
RUN wget -nv -O "/tmp/shellcheck-${SHELLCHECK}.linux.$(uname -m).tar.xz" "https://github.com/koalaman/shellcheck/releases/download/${SHELLCHECK}/shellcheck-${SHELLCHECK}.linux.$(uname -m).tar.xz"
RUN tar -xJf "/tmp/shellcheck-${SHELLCHECK}.linux.$(uname -m).tar.xz" -C /tmp
RUN mv /tmp/shellcheck-${SHELLCHECK}/shellcheck /usr/bin

# trivy
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN wget -nv -O "/tmp/trivy_${TRIVY}_Linux-64bit.deb" "https://github.com/aquasecurity/trivy/releases/download/v${TRIVY}/trivy_${TRIVY}_Linux-64bit.deb" && \
    apt-get -y install --no-install-recommends -f "/tmp/trivy_${TRIVY}_Linux-64bit.deb" && \
    rm "/tmp/trivy_${TRIVY}_Linux-64bit.deb"

RUN go install mvdan.cc/gofumpt@${GOFUMPT} && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO} && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC} && \
    go install github.com/envoyproxy/protoc-gen-validate@${PROTOC_GEN_VALIDATE} && \
    go install sigs.k8s.io/kind@${KIND} && \
    go clean -modcache && \
    PATH=${PATH}:$(go env GOPATH)/bin && \
    export PATH

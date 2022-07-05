#!/bin/bash

set -euo pipefail

buf breaking --against ".git#branch=main"
buf lint

gofumpt -l -w .
golangci-lint run

go mod tidy
podman build . -t jacodelg/playground:latest

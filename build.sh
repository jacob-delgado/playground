#!/bin/bash

set -euo pipefail

buf lint
go mod tidy
gofumpt -l -w .
golangci-lint run
podman build . -t jacodelg/playground:latest

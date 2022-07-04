#!/bin/bash

set -euo pipefail

gofumpt -l -w .
golangci-lint run
podman build . -t jacodelg/playground:latest

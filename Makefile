all: inventory storefront

.PHONY: build-tools
build-tools:
	docker build --platform linux/amd64 -f build-tools/Dockerfile -t jdelgad/build-tools:latest .

.PHONY: clean
clean:
	rm -rf gen
	rm inventory || true
	rm storefront || true

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: gen
gen:
	buf generate

.PHONY: hadolint
hadolint:
	hadolint build-tools/Dockerfile
	hadolint cmd/inventory/Dockerfile
	hadolint cmd/storefront/Dockerfile

.PHONY: images
images: clean lint test
	docker build --platform linux/amd64 -f cmd/inventory/Dockerfile -t jdelgad/inventory:latest .
	docker build --platform linux/amd64 -f cmd/storefront/Dockerfile -t jdelgad/storefront:latest .

.PHONY: inventory
inventory: clean lint test
	GOOS=linux go build ./cmd/inventory

.PHONY: lint
lint: gen fmt tidy
	golangci-lint run
	buf lint

.PHONY: push-build-tool
push-build-tools:
	docker push jdelgad/build-tools:latest

.PHONY: shellcheck
shellcheck:
	shellcheck setup.sh
	shellcheck grpc-call.sh

.PHONY: tidy
	go mod download
	go mod tidy

.PHONY: storefront
storefront: clean lint test
	GOOS=linux go build ./cmd/storefront

.PHONY: test
test: clean gen
	go test ./... -race

.PHONY: clean-images
clean-images:
	docker rmi jdelgad/inventory:latest
	docker rmi jdelgad/storefront:latest

.PHONY: prune
prune:
	docker system prune --all

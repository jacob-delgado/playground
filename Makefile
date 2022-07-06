all: inventory

.PHONY: clean
clean:
	rm -rf gen
	rm inventory || true

.PHONY: gen
gen:
	buf generate

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: tidy
	go mod download
	go mod tidy

.PHONY: lint
lint: gen fmt tidy
	golangci-lint run
	buf lint

.PHONY: test
test: clean gen
	go test ./... -race

.PHONY: inventory
inventory: clean lint test
	GOOS=linux go build ./cmd/inventory

.PHONY: image
image: clean lint test
	docker build . -t jacodelg/inventory:latest

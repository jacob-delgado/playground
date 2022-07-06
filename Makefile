all: playground

.PHONY: clean
clean:
	rm -rf gen
	rm playground || true

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

.PHONY: playground
playground: clean lint test
	GOOS=linux go build ./cmd/playground

.PHONY: image
image: clean lint test
	docker build . -t jacodelg/playground:latest

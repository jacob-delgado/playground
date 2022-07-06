all: inventory storefront

.PHONY: clean
clean:
	rm -rf gen
	rm inventory || true
	rm storefront || true

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

.PHONY: storefront
storefront: clean lint test
	GOOS=linux go build ./cmd/storefront

.PHONY: images
images: clean lint test
	docker build -f cmd/inventory/Dockerfile -t jacodelg/inventory:latest .
	docker build -f cmd/storefront/Dockerfile -t jacodelg/storefront:latest .

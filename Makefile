.PHONY: build vet lint test check

build:
	go build ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test -race -cover ./...

check: build vet lint test

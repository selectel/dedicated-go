default: tests

tests: golangci-lint test

test:
	go test ./...

golangci-lint:
	docker run --rm -v $(PWD):/app:ro -w /app golangci/golangci-lint:v2.12.2 golangci-lint run

.PHONY: tests unittest golangci-lint
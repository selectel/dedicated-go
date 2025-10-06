default: tests

tests: golangci-lint test

test:
	go test ./...

golangci-lint:
	@sh -c "'$(CURDIR)/scripts/golangci_lint_check.sh'"

.PHONY: tests unittest golangci-lint
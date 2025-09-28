GOLANGCILINT=golangci-lint
MODERNIZE=go tool modernize

.PHONY: all
all:

.PHONY: lint
lint: lint-go

.PHONY: lint-go
lint-go:
	go vet ./...
ifneq ($(CI),1)
	$(GOLANGCILINT) run
endif

.PHONY: lint-go-modernize
lint-go-modernize:
	$(MODERNIZE) -test ./...

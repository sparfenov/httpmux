REFLEX ?= github.com/cespare/reflex
GOLANGCI_LINT ?= $(GOPATH)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.30.0
SEP			 ?= "========================================================"
##
.PHONY: env

define ENV_SAMPLE
SERVER_ADDR=:8080
IS_DEBUG=true
endef

export ENV_SAMPLE
env:
	@if [ ! -f ".env" ];\
        then echo "$$ENV_SAMPLE" > .env;\
        echo ".env created";\
    else\
        echo ".env already exists";\
    fi

################################################################################################################
.PHONY: dev
dev:
	go run $(REFLEX) -R "\\.idea|vendor" -r "\\.go" -s -- sh -c "go run --race ./cmd/app/..."

################################################################################################################
.PHONY: gofmt
gofmt:
	gofmt -w cmd/* internal/* pkg/*

.PHONY: lint
lint:
	@make gofmt

	$(GOLANGCI_LINT) run --config .golangci.yml ./...

.PHONY: install-golangci-lint
install-golangci-lint:
	mkdir -p $(GOPATH)/bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_LINT_VERSION)/install.sh \
		| sed -e '/install -d/d' \
		| sh -s -- -b $(GOPATH)/bin $(GOLANGCI_LINT_VERSION)

.PHONY: test
test:
	$(call _info, $(SEP))
	$(call _info,"RUN TESTS")
	$(call _info, $(SEP))
	GO111MODULE=on go test -coverprofile coverage.out ./internal/...
	go tool cover -func coverage.out
################################################################################################################

define _info
	$(call _echoColor,$1,6)
endef

define _hint
	$(call _echoColor,$1,8)
endef

define _succ
	$(call _echoColor,$1,2)
endef

define _warn
	$(call _echoColor,$1,3)
endef

define _mega
	$(call _echoColor,$1,13)
endef

define _error
	$(call _echoColor,$1,1)
endef

define _echoColor
	@tput setaf $2
	@echo $1
	@tput sgr0
endef

################################################################################################################

aws-deploy:
	make build
	serverless deploy --verbose

aws-redeploy:
	serverless remove
	make aws-deploy
# docker run --rm -it amazon/aws-cli command

PROJECT_BIN ?= $(shell pwd)/bin

.PHONY: precommit
precommit: modtidy modupdate gofumpt lint test govuln

.PHONY: modtidy
modtidy:
	go mod tidy

.PHONY: modupdate
modupdate:
	go get -u ./...

.PHONY: gofumpt
gofumpt: install-gofumpt
	$(GOFUMPT_BIN) -l -w -extra .

.PHONY: lint
lint: install-lint
	$(LINT_BIN) run -v

.PHONY: test
test:
	go test -v ./...

.PHONY: govuln
govuln: install-govuln
	$(GOVULN_BIN) ./...

GOFUMPT_BIN = $(PROJECT_BIN)/gofumpt
.PHONY: install-gofumpt
install-gofumpt:
	$(call go-install,$(GOFUMPT_BIN),mvdan.cc/gofumpt)

LINT_BIN = $(PROJECT_BIN)/golangci-lint
.PHONY: install-lint
install-lint:
	$(call go-install,$(LINT_BIN),github.com/golangci/golangci-lint/cmd/golangci-lint)

GOVULN_BIN = $(PROJECT_BIN)/govulncheck
.PHONY: install-govuln
install-govuln:
	$(call go-install,$(GOVULN_BIN),golang.org/x/vuln/cmd/govulncheck)

define go-install
	@if [ ! -f $1 ] ; then \
		echo "Installing $1 to the $(PROJECT_BIN) ..."; \
		GOBIN=$(PROJECT_BIN) go install $2@latest; \
	fi
endef

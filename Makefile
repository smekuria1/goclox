.DEFAULT_GOAL := build-win
BIN_FILE=goclox.exe

MODULE   = $(shell $(GO) list -m)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat .version 2> /dev/null || echo v0)
PKGS     = $(or $(PKG),$(shell $(GO) list ./...))
BIN      = bin

GO      = go
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell if [ "$$(tput colors 2> /dev/null || echo 0)" -ge 8 ]; then printf "\033[34;1m▶\033[0m"; else printf "▶"; fi)

GENERATED = # List of generated files
# Tools

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)…)
	$Q env GOBIN=$(abspath $(BIN)) $(GO) install $(PACKAGE)

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports@latest

REVIVE = $(BIN)/revive
$(BIN)/revive: PACKAGE=github.com/mgechev/revive@v1.2.4

GOCOV = $(BIN)/gocov
$(BIN)/gocov: PACKAGE=github.com/axw/gocov/gocov@latest

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml@latest

GOTESTSUM = $(BIN)/gotestsum
$(BIN)/gotestsum: PACKAGE=gotest.tools/gotestsum@latest

RUNMODE = repl

build-win:
	@echo 'Building binary in ./bin/${BIN_FILE}'
	@go build -o "./bin/${BIN_FILE}"
clean:
	go clean
	rm --force "cp.out"
	rm --force nohup.out
	rm -r --force ./bin
test:
	go test ./...
check:
	go test
cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out
run:
	./bin/"${BIN_FILE}" -file test.clox
run-exT:
	./bin/"${BIN_FILE}" -file test.clox -debugT

run-exC:
	./bin/"${BIN_FILE}" -file test.clox -debugC

lint: | $(REVIVE) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(REVIVE) -formatter friendly -set_exit_status ./...

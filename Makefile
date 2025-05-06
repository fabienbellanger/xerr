.PHONY: update \
	update-all \
	format \
	lint \
	test \
	test-verbose \
	test-tparse \
	bench \
	clean \
	doc \
	help \
	test-cover-count \
	cover-count \
	test-cover-atomic \
	cover-atomic \
	html-cover-count \
	html-cover-atomic \
	run-cover-count \
	run-cover-atomic \
	view-cover-count \
	view-cover-atomic

.DEFAULT_GOAL=help

# Read: https://kodfabrik.com/journal/a-good-makefile-for-go

# Go parameters
CURRENT_PATH=$(shell pwd)
GO_CMD=go
GO_RUN=$(GO_CMD) run
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_MOD=$(GO_CMD) mod
GO_TOOL=$(GO_CMD) tool
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt
GODOC=godoc

## update: Update modules
update:
	$(GO_GET) -u ./... && $(GO_MOD) tidy

## update-all: Update all modules
update-all:
	$(GO_GET) -u ./... all && $(GO_MOD) tidy

## format: Run go fmt
format:
	$(GO_FMT) ./...

## lint: Run go vet
lint: format
	$(GO_VET) ./...

## test: Run test
test:
	$(GO_TEST) -cover ./...

## test-verbose: Run tests
test-verbose:
	$(GO_TEST) -cover -v ./...

## test-tparse: Run tests with tparse
test-tparse:
	go test -cover -json ./... | tparse -trimpath -all

test-cover-count:
	$(GO_TEST) -covermode=count -coverprofile=cover-count.out ./...

test-cover-atomic:
	$(GO_TEST) -covermode=atomic -coverprofile=cover-atomic.out ./...

cover-count:
	$(GO_TOOL) cover -func=cover-count.out

cover-atomic:
	$(GO_TOOL) cover -func=cover-atomic.out

html-cover-count:
	$(GO_TOOL) cover -html=cover-count.out

html-cover-atomic:
	$(GO_TOOL) cover -html=cover-atomic.out

run-cover-count: test-cover-count cover-count
	rm cover-count.out
run-cover-atomic: test-cover-atomic cover-atomic
	rm cover-atomic.out
view-cover-count: test-cover-count html-cover-count
	rm cover-count.out
view-cover-atomic: test-cover-atomic html-cover-atomic
	rm cover-atomic.out

## bench: Run benchmarks
bench:
	$(GO_TEST) -benchmem -bench=. ./...

## clean: Clean files
clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

## doc: Launch godoc on port 9898
doc:
	$(GODOC) -http :9898

help: Makefile
	@echo
	@echo "Choose a command run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo

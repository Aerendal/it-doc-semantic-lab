BINARY     := itdlab
BUILD_DIR  := bin
CMD_PATH   := ./cmd/itdlab
DB_PATH    := db/semantic_index.sqlite
LOG_PATH   := runs/events.jsonl

.PHONY: all build test lint clean run help

all: build

## build: compile the CLI binary into bin/
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD_PATH)

## test: run all tests
test:
	go test ./... -v

## test-short: run tests without integration tests
test-short:
	go test ./... -short -v

## lint: run go vet (add golangci-lint when available)
lint:
	go vet ./...

## tidy: tidy go.mod and go.sum
tidy:
	go mod tidy

## clean: remove build artifacts and local runtime state
clean:
	rm -rf $(BUILD_DIR)
	rm -f $(DB_PATH) $(LOG_PATH)

## run-ingest: ingest sources from default source directory
run-ingest: build
	./$(BUILD_DIR)/$(BINARY) ingest run

## run-normalize: preview normalization
run-normalize: build
	./$(BUILD_DIR)/$(BINARY) normalize preview

## db-init: initialise a fresh SQLite database (via schema_v1.sql)
db-init:
	@mkdir -p db runs
	sqlite3 $(DB_PATH) < db/schema_v1.sql
	@echo "Database initialised at $(DB_PATH)"

## db-schema: show current schema version
db-schema:
	sqlite3 $(DB_PATH) "SELECT * FROM schema_version;"

## help: show this help
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'

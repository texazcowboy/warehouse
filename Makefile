API_SOURCES := $(shell find ./cmd/warehouse-api -type f -name "*.go" -not -name "*_test.go" )
API_BINARY := ./bin/warehouse-api
MIGRATION_SOURCES := $(shell find ./cmd/warehouse-migration -type f -name "*.go" -not -name "*_test.go" )
MIGRATION_BINARY := ./bin/warehouse-migration
CONFIG_LOCATION := ./config.yaml

.PHONY: run-api
run-api: build-api
	CONFIG=$(CONFIG_LOCATION) ./bin/warehouse-api

.PHONY: build-api
# just an alias for api binary path
build-api: $(API_BINARY)

$(API_BINARY): $(API_SOURCES)
	go build -o $(API_BINARY) -v ./cmd/warehouse-api

.PHONY: test-api
test-api:
	go test ./cmd/warehouse-api/... -v

.PHONY: run-migration
run-migration: build-migration
	CONFIG=$(CONFIG_LOCATION) SOURCE=file://cmd/warehouse-migration/migrations DIRECTION=up ./bin/warehouse-migration

.PHONY: build-migration
# just an alias for migration binary path
build-migration: $(MIGRATION_BINARY)

$(MIGRATION_BINARY): $(MIGRATION_SOURCES)
	go build -o $(MIGRATION_BINARY) -v ./cmd/warehouse-migration/

.PHONY: test-migration
test-migration:
	go test ./cmd/warehouse-migration/... -v

.PHONY: clean
clean:
	go clean
	rm -r ./bin

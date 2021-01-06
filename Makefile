CONFIG_LOCATION=./config.yaml

run-api: build-api
	./bin/warehouse-api -config=$(CONFIG_LOCATION)

build-api: test-api
	go build -o ./bin/warehouse-api -v ./cmd/warehouse-api

test-api:
	go test ./cmd/warehouse-api/... -v

run-migration: build-migration
	./bin/warehouse-migration -config=$(CONFIG_LOCATION) -src=file://cmd/warehouse-migration/migrations

build-migration: test-migration
	go build -o ./bin/warehouse-migration -v ./cmd/warehouse-migration/

test-migration:
	go test ./cmd/warehouse-migration/... -v

clean:
	go clean
	rm -r ./bin

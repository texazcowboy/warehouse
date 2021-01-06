run-api: build-api
	./bin/warehouse-api -config=config.yaml

test-api:
	go test ./cmd/warehouse-api/... -v

build-api:
	go build -o ./bin/warehouse-api -v ./cmd/warehouse-api

run-migration: build-migration
	./bin/warehouse-migration -config=config.yaml -src=file://cmd/warehouse-migration/migrations

test-migration:
	go test ./cmd/warehouse-migration/... -v

build-migration:
	go build -o ./bin/warehouse-migration -v ./cmd/warehouse-migration/

clean:
	go clean
	rm -r ./bin

all: test build migrate run

test_unit:
	go test ./... -v

test_integration:
	INTEGRATION=1 go test -v ./integration_tests/...

test: test_unit test_integration

build:
	docker build --tag $(USER)/helsinki-guide:$(TAG) .

migrate:
	migrate -database "${DATABASE_URL}" -path internal/infrastructure/migrations up

run:
	docker pull andreyad/helsinki-guide
	docker run \
	--env DEBUG=1 \
	--env DATABASE_URL="${DATABASE_URL}" \
	--env BOT_TOKEN="${BOT_TOKEN}" \
	--network host \
	--log-opt tag=hguide \
	--name helsinki-guide \
	andreyad/helsinki-guide

.NOTPARALLEL:

.PHONY: all test_unit test_integration test build migrate run
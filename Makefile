all: test build migrate run

test_unit:
	go test ./... -v

test_integration:
	INTEGRATION=1 go test ./integration_tests/...

test: test_unit test_integration

build:
	docker build --tag andreyad/helsinki-guide:latest .

migrate:
	migrate -database "${DatabaseURL}" -path internal/infrastructure/migrations up

run:
	docker pull andreyad/helsinki-guide
	docker run \
	--env Debug=1 \
	--env DatabaseURL="${DatabaseURL}" \
	--env BotAPIToken="${BotAPIToken}" \
	--network host \
	--log-opt tag=hguide \
	helsinki-guide

.NOTPARALLEL:

.PHONY: test_unit test_integration test build migrate run
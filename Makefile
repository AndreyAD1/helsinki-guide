build:
	docker build .

migrate:
	migrate -database ${DatabaseURL} -path internal/infrastructure/migrations up

run:
	docker run helsinki-guide

.NOTPARALLEL:

.PHONY: build migrate run
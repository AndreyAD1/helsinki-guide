build:
	docker build --tag helsinki-guide:latest .

migrate:
	migrate -database ${DatabaseURL} -path internal/infrastructure/migrations up

run:
	docker run --env DatabaseURL="${DatabaseURL}" --env BotAPIToken="${BotAPIToken}" helsinki-guide

.NOTPARALLEL:

.PHONY: build migrate run
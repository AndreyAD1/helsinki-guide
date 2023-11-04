build:
	docker build .

migrate:
	python $(CURDIR)/scripts/migrate.py

run:
	docker run helsinki-guide

.NOTPARALLEL:

.PHONY: build migrate run
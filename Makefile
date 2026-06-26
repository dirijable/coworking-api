include .env
export

PROJ_ROOT=$(shell pwd)

cwa-run:
	@go mod tidy && \
	go run cmd/coworking-api/main.go

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "value 'seq' is required"; \
		exit 1; \
	fi; \
	docker compose run --rm migrations \
		create \
		-ext sql \
		-dir ./migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
  		echo "value 'action' is required"; \
  		exit 1; \
  	fi; \
	docker compose run --rm migrations \
		-path ./migrations \
		-database "postgres://${PG_USER}:${PG_PASSWORD}@coworking-db:5432/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}" \
		"$(action)"

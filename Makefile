.PHONY: build run migrate lint format

# Build all Docker services
build:
	docker-compose build

# Run the Docker services
run:
	docker-compose up

# Apply database migrations using Flyway
migrate:
	docker-compose run flyway migrate

# Run linting using golangci-lint in the backend service
lint:
	docker-compose run backend golangci-lint run

# Format the Go code using goimports in the backend service
format:
	docker-compose run backend goimports -w ./internal

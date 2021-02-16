SERVICE = parenthesis_service

run: lint test run-prom
	@echo "Run service"
	@docker-compose up parenthesis_service

test:
	@go test -cover ./...

lint:
	@golangci-lint run -c .golangci.yml ./...

run-prom:
	@docker-compose up -d prometheus

down:
	@echo "Turn off"
	@docker-compose down

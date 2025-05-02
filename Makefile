# Set up commands
init-vscode:
	cd src/server && cp .env.sample .env && cd .. && cd .. && \
	cd .vscode && cp settings.json.sample settings.json

# Set up database
set-up-db:
	docker compose down && \
	docker compose up -d server-postgres
	cd src/server/ && go run tests/commands/clear_database.go

# Test commands
test:
	make set-up-db && \
	make test-go

test-go:
	go test -p 1 -count=1 ./src/server/tests/...

# Runners
run:
	cd src/server && go run main.go

# Swagger
swag-server:
	cd src/server/api && swag init -g server.go --instanceName server --parseDependency --parseDepth 1

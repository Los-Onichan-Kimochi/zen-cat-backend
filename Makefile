# Set up commands
init-vscode:
	cd src/server && cp .env.sample .env && cd .. && cd .. && \
	cd .vscode && cp settings.json.sample settings.json

# Set up database (run with sudo)
set-up-db:
	docker compose down && \
	docker compose up astro-cat-postgres -d --wait
	cd src/server/ && go run tests/set_dummy_data.go

# Set up testing enviroment for python
setup-test:
	cd src/server/tests && \
	python3 -m venv libs && \
	. libs/bin/activate && \
	pip install -r requirements.txt && \
	deactivate

# Runners
run:
	cd src/server && go run main.go

# Test commands
test:
	cd src/server/ && \
	go run tests/set_dummy_data.go && \
	cd .. && cd .. && \
	cd src/server/tests && \
	. libs/bin/activate && \
	python3 -m pytest

# Swagger documentation
swag-docs:
	cd src/server/api && swag init -g server.go --instanceName server --parseDependency --parseDepth 1

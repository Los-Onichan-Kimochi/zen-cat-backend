# Set up commands
init-vscode:
	cd src/server && cp .env.sample .env && cd .. && cd .. && \
	cd .vscode && cp settings.json.sample settings.json

# Set up AWS S3 credentials
set-aws-credentials:
	@./scripts/setup-aws-s3.sh

# Set up database
set-up-db:
	docker compose down && \
	docker compose up astro-cat-postgres -d --wait
	cd src/server/ && go run tests/commands/clear_database.go

# Test commands
test:
	make set-up-db && \
	make test-go

test-go:
	go test -p 1 -count=1 ./src/server/tests/...

test-go-sum:
	gotestsum --format pkgname  -- -p 1 -count=1 ./src/server/tests/...

# Runners
run:
	cd src/server && go run main.go

# Swagger documentation
swag-docs:
	cd src/server/api && swag init -g server.go --instanceName server --parseDependency --parseDepth 1

# S3 test
s3-test:
	cd src/server/tests/s3_test && go run s3.go

# Railway deployment commands
railway-build:
	@./scripts/build-for-railway.sh

railway-deploy:
	@echo "🚀 Deploying to Railway..."
	@echo "📝 Make sure you have Railway CLI installed and are logged in"
	@echo "📝 Run: railway up"

railway-logs:
	@echo "📋 Viewing Railway logs..."
	@echo "📝 Run: railway logs"

railway-vars:
	@echo "🔧 Railway environment variables:"
	@echo "📝 Run: railway variables"

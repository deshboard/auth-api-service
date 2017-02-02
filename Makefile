GLIDE:=$(shell if which glide > /dev/null 2>&1; then echo "glide"; fi)

# Setup environment
setup: build
	mkdir -p var/
	docker-compose up -d db
	make migrate

# Build the service and test containers
build:
ifeq ($(FORCE), true)
	@docker-compose build --force-rm db service test
else
	@docker-compose build db service test
endif

# Run database migrations
migrate:
	@docker-compose --file docker-compose.yml --file docker-compose.util.yml run --rm db_check
	@docker-compose --file docker-compose.yml --file docker-compose.util.yml run --rm service.model.user.migration update

# Start the environment
start:
	@docker-compose up -d

# Stop the environment
stop:
ifeq ($(FORCE), true)
	@docker-compose kill
else
	@docker-compose stop
endif

# Clean environment
clean: stop
	@rm -rf vendor/ var/
	@docker-compose rm --force
	@go clean

# Run test suite
test:
	@docker-compose exec db mysql -e "DELETE FROM service_model_user.users"
	@docker-compose run --rm test

# Install dependencies locally, optionally using go get
install:
ifdef GLIDE
	@$(GLIDE) install
else ifeq ($(FORCE), true)
	@go get
else
	@echo "Glide is necessary for installing project dependencies: http://glide.sh/ Run this command with FORCE=true to fall back to go get" 1>&2 && exit 1
endif

# Check that all source files follow the Coding Style
cs:
	@gofmt -l . | read something && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true

# Fix Coding Standard violations
csfix:
	@gofmt -l -w -s .

.PHONY: setup build migrate start stop clean test install cs csfix

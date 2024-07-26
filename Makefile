# Define service names
SERVICES := config-loader random-modifier service-6b631a97 service-9fd2041b service-e079e8d7
# Docker Compose file
DOCKER_COMPOSE_FILE := docker-compose.yml

# Default target: build and run everything
.PHONY: all
all: build up

# Build all Docker images
.PHONY: build
build:
	@echo "Building Docker images..."
	$(foreach service, $(SERVICES), docker compose -f $(DOCKER_COMPOSE_FILE) build $(service);)

# Start all services
.PHONY: up
up:
	@echo "Starting services..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d
	@docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Stop all services
.PHONY: down
down:
	@echo "Stopping services..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up: stop services and remove containers, networks, volumes, and images
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) down --rmi all -v --remove-orphans

# Tail logs for all services
.PHONY: logs
logs:
	@docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Restart services
.PHONY: restart
restart: down up

# Display help
.PHONY: help
help:
	@echo "Makefile for managing Docker containers"
	@echo
	@echo "Usage:"
	@echo "  make          Build and start all services"
	@echo "  make build    Build all Docker images"
	@echo "  make up       Start all services"
	@echo "  make down     Stop all services"
	@echo "  make clean    Stop services and remove containers, networks, volumes, and images"
	@echo "  make logs     Tail logs for all services"
	@echo "  make restart  Restart all services"
	@echo "  make help     Display this help message"
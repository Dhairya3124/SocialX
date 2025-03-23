# Makefile for Goose migrations
#include .envrc
include .envrc
# Variables
GOOSE_CMD=goose
DB_DIR=./cmd/migrate/migrations/schema

# Targets
.PHONY: create-migration up down status

create-migration:
	@read -p "Enter migration name: " name; \
	$(GOOSE_CMD) -s -dir $(DB_DIR) create $$name sql

up:
	$(GOOSE_CMD) postgres -dir $(DB_DIR) $(DB_ADDR) up

down:
	$(GOOSE_CMD) postgres -dir $(DB_DIR) $(DB_ADDR) down


fmt:
	go fmt ./...

.PHONY: gen-docs
gen-docs:
	swag init -g ./api/main.go -d cmd,internal && swag fmt
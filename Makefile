GO=go
SQLC=sqlc
GOOSE=goose
SWAG=swag
API_DIR=api

.DEFAULT_GOAL := help
.PHONY: help dev build docs

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\) \(.*\)##\(.*\)/\1\3/p'

dev: ## Run the development server
	$(MAKE) -C $(API_DIR) sqlc-generate
	$(MAKE) -C $(API_DIR) goose-up
	$(MAKE) -C $(API_DIR) go-populate
	$(MAKE) -C $(API_DIR) swag-fmt
	$(MAKE) -C $(API_DIR) swag-init
	$(MAKE) -C $(API_DIR) go-run

down: ## Tear down the development database
	$(MAKE) -C $(API_DIR) goose-down

build: ## Build the whole project
	$(MAKE) -C $(API_DIR) sqlc-generate
	$(MAKE) -C $(API_DIR) swag-fmt
	$(MAKE) -C $(API_DIR) swag-init
	$(MAKE) -C $(API_DIR) go-build

docs: ## Build the documentation
	$(MAKE) -C $(API_DIR) swag-fmt
	$(MAKE) -C $(API_DIR) swag-init
GO=go
GOOSE=${GO} tool github.com/pressly/goose/v3/cmd/goose
SWAG=${GO} tool github.com/swaggo/swag/cmd/swag

.PHONY: sqlc-generate goose-up goose-down swag-fmt swag-init go-run go-build

# Run migrations
goose-up:
	$(GOOSE) up

# Undo the latest migration
goose-down:
	$(GOOSE) down

# Run the server
go-run:
	./scripts/start.sh

# Build the project
go-build:
	$(GO) generate
	$(GO) build -o eigakanban
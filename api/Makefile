GO=go
SQLC=sqlc
GOOSE=goose
SWAG=swag

.PHONY: sqlc-generate goose-up goose-down swag-fmt swag-init go-run go-build

# Generate SQL code
sqlc-generate:
	$(SQLC) generate

# Run migrations
goose-up:
	$(GOOSE) up

# Undo the latest migration
goose-down:
	$(GOOSE) down

# Format code and Swagger docs
swag-fmt:
	$(SWAG) fmt && $(GO) fmt

# Initialize Swagger docs
swag-init:
	$(SWAG) init

# Populate the database with dummy values
go-populate:
	$(GO) run scripts/insert_dummy_data.go

# Run the server
go-run:
	$(GO) run main.go

# Build the project
go-build:
	$(GO) build -o eigakanban
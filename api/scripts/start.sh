# Generate code
go generate

# Run database migrations
echo "Running migrations..."
go tool github.com/pressly/goose/v3/cmd/goose up

# Insert dummy data
echo "Inserting dummy data..."
go run scripts/insert_dummy_data.go

# Start the application
echo "Starting eigakanban..."
go tool github.com/air-verse/air

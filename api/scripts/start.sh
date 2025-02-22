# Run database migrations
echo "Running migrations..."
goose up

# Insert dummy data
echo "Inserting dummy data..."
go run scripts/insert_dummy_data.go

# Start the application
echo "Starting eigakanban..."
air

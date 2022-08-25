migrateup:
	migrate -path pkg/migrations -database "postgresql://postgres:demo_password@localhost:5432/Bank?sslmode=disable" -verbose up
migratedown:
	migrate -path pkg/migrations -database "postgresql://postgres:demo_password@localhost:5432/Bank?sslmode=disable" -verbose down
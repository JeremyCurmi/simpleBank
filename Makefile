migrate_up:
	migrate -path pkg/migrations -database "postgresql://postgres:demo_password@localhost:5432/Bank?sslmode=disable" -verbose up
migrate_down:
	yes | migrate -path pkg/migrations -database "postgresql://postgres:demo_password@localhost:5432/Bank?sslmode=disable" -verbose down
docker_up:
	docker compose up -d --build
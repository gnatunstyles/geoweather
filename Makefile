.PHONY: migrate-up
migrate-up:
			docker exec -i db psql -U postgres < db/migrate.sql

.PHONY: migrate-down
migrate-down:
			docker exec -i db psql -U postgres < db/migrate_down.sql

.PHONY: run
run:
	go run cmd/main.go

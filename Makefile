include .env

migrate-up:
	migrate -path db/migration -database "postgresql://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

tidy:
	go mod tidy
include .env

migrate-up:
	migrate -path db/migration -database ${DB_URL} -verbose up

migrate-down:
	migrate -path db/migration -database ${DB_URL} -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

tidy:
	go mod tidy
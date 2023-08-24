include .env

clean:
	go clean -testcache

migrate-up:
	migrate -path db/migration -database ${DB_URL} -verbose up

migrate-down:
	migrate -path db/migration -database ${DB_URL} -verbose down

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/drmanalo/simplebank/db/sqlc Store

server:
	go run main.go

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

tidy:
	go mod tidy

.PHONY: clean migrate-up migrate-down mock server sqlc test tidy
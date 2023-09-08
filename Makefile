include .env

clean:
	go clean -testcache

migrate-up:
	migrate -path db/migration -database ${DB_URL} -verbose up

migrate-up1:
	migrate -path db/migration -database ${DB_URL} -verbose up 1

migrate-down:
	migrate -path db/migration -database ${DB_URL} -verbose down

migrate-down1:
	migrate -path db/migration -database ${DB_URL} -verbose down 1

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
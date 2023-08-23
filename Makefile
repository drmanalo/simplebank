include .env

clean:
	go clean -testcache

migrate-up:
	migrate -path db/migration -database ${DB_URL} -verbose up

migrate-down:
	migrate -path db/migration -database ${DB_URL} -verbose down

server:
	go run main.go

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

tidy:
	go mod tidy
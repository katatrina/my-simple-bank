postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=simple_bank -d postgres:12-alpine

create-server:
	docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://postgres:secret@postgres12:5432/simple_bank?sslmode=disable" simplebank:latest

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v $(CURDIR):/src -w /src sqlc/sqlc generate

server:
	go run main.go

## Run unit tests for all packages
test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/katatrina/my-simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server
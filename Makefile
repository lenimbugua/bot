DB_SOURCE=postgresql://root:password@localhost:5432/bot?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14.5-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root bot

dropdb:
	docker exec -it postgres dropdb bot

initschema: 
	migrate create -ext sql -dir db/migration init_schema

prepopulate: 
	migrate create -ext sql -dir db/migration prepopulate_db

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/lenimbugua/bot/db/sqlc Store

.PHONY: postgres createdb dropdb initschema migrateup prepopulate migratedown migrateup1 migratedown1 sqlc test server mock

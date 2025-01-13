# golang-crud

Create db migration

```
docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
brew install golang-migrate
migrate create -ext sql -dir db/migration -seq init_schema
make migrateup
```

Generate SQLC

```
sqlc init
sqlc generate
```

Run

```
go run main.go
```

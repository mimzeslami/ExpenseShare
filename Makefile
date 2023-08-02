createdb:
	docker exec -it postgres12 createdb --username=root --owner=root expense_share
dropdb:
	docker exec -it postgres12 dropdb expense_share
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/expense_share?sslmode=disable" -verbose up 
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/expense_share?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test: 
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mimzeslami/expense_share/db/sqlc Store

.PHONE: createdb dropdb migrateup migratedown sqlc test server mock
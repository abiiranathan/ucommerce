DATABASE=postgres://postgres:ucommerce@127.0.0.1:5432/ucommerce?sslmode=disable
MIGRATIONS=file://postgres/migrations

migrate:
	migrate -source $(MIGRATIONS) -database $(DATABASE) up

down:
	migrate -source $(MIGRATIONS) -database $(DATABASE) down

drop:
	migrate -source $(MIGRATIONS) -database $(DATABASE) drop

migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir postgres/migrations $$name

sqlc:
	sqlc generate

psql:
	psql -U nabiizy sqlc
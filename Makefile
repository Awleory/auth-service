test:
	go test -v -cover ./internal/storage/psql
migrate:
	migrate -path ./schema -database 'postgres://postgres:so2037456va@0.0.0.0:5432/postgres?sslmode=disable' up
migrate-down:
	migrate -path ./schema -database 'postgres://postgres:so2037456va@0.0.0.0:5432/postgres?sslmode=disable' down
build:
	docker-compose build
run:
	docker-compose up
create-m:
	migrate create -ext sql -dir ./schema -seq init
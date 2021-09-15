.PHONY: migrate

migrate:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/baby_guess?sslmode=disable" -verbose up

run:
	HTTP_SERVER_ADDR=":8080" go run .

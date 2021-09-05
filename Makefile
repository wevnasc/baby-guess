migrate:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5432/baby_guess?sslmode=disable" -verbose up                   ~/dev/wevnasc/baby-guess

.PHONY: migrate

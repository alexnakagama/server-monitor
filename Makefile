include .env
export

migrate-up:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)" down 1

migrate-test-up:
	migrate -path db/migrations -database "$(TEST_DATABASE_URL)" up

migrate-test-down:
	migrate -path db/migrations -database "$(TEST_DATABASE_URL)" down 1

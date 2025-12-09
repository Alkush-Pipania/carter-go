# Load from environment, fallback to .env file
include .env
export

# Run migrations up
migrate-up:
	migrate -path migration -database "$(DB_URL)" up

# Rollback last migration
migrate-down:
	migrate -path migration -database "$(DB_URL)" down 1

# Force a specific version (use if migration gets stuck)
migrate-force:
	migrate -path migration -database "$(DB_URL)" force $(VERSION)

# Create a new migration
migrate-create:
	migrate create -ext sql -dir migration -seq $(NAME)
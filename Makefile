include .env

POSTGRES_DSN = "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

migration-create: $(MIGRATE) ## Create a set of up/down migrations with a specific name
	@ read -p "Please provide name for the migration: " Name; \
	migrate create -ext sql -dir databases -seq $${Name}

migration-up: $(MIGRATE) ## Apply all (or N up) migrations
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate -database $(POSTGRES_DSN) -path=databases up $(N)

migration-down: $(MIGRATE) ## Apply all (or N down) migrations
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate -database $(POSTGRES_DSN) -path=databases down $(N)
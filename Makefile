include .env

MIGRATIONS_PATH=./migrations

.PHONY: migration
migration:
	@migrate create -seq -ext=.sql -dir=$(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

dbu:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_DSN) -verbose up $(filter-out $@,$(MAKECMDGOALS))

dbd:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_DSN) -verbose down $(filter-out $@,$(MAKECMDGOALS))

ideps:
	@cd web/ && pnpm i

dev:
	@docker compose up -d 
	@sleep 2
	@air &
	@sleep 2
	@cd web && pnpm dev &
	@sleep 2
	@caddy run --config Caddyfile.dev &

stop:
	@docker-compose down
	-@pkill -x "air"
	-@pkill -f "vite"
	-@caddy stop

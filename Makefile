DB_ADDR=postgres://postgres:postgres@localhost:5432/moota?sslmode=disable
MIGRATIONS_PATH=./migrations

.PHONY: migration
migration:
	@migrate create -seq -ext=.sql -dir=$(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

dbu:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose up $(filter-out $@,$(MAKECMDGOALS))

dbd:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose down $(filter-out $@,$(MAKECMDGOALS))

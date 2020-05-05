export HTTP_HOST := localhost
export HTTP_ADDR := :8198
export LOG_LEVEL := debug

export DATABASE_HOST := localhost
export DATABASE_NAME := bamboo
export DATABASE_USER := bamboo
export DATABASE_PASSWORD := bamboo
export DATABASE_SSLMODE := disable

export MIGRATIONS_FILES := /migrations


export GOPROXY=direct

.PHONY: configs.env
configs.env:
	cp config/ragger_example.env config/ragger.env
	cp config/database_example.env config/database.env

.PHONY: test
test:
	go test -v -race ./...

.PHONY: run
run:test
	go run ./cmd/main.go

.PHONY: swagger
swagger:
	$(HOME)/go/bin/swag init -g ./cmd/main.go

export POSTGRESQL_URL='postgres://bamboo:bamboo@172.17.0.2:31770/bamboo?sslmode=disable'

.PHONY: migrate
migrate:
	migrate -database ${POSTGRESQL_URL} -path migrations up
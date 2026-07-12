.PHONY: swag certs env migrate-up migrate-down migrate-status migrate-reset migrate-create

include .env
export

dsn = "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

swag:
	swag init -g cmd/web/main.go --output cmd/web/docs

certs:
	mkdir certs
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -pubout -out public.pem
	mv private.pem public.pem certs/

env:
	cp .env.example .env

migrate-up:
	goose -dir migrations postgres $(dsn) up

migrate-down:
	goose -dir migrations postgres $(dsn) down

migrate-status:
	goose -dir migrations postgres $(dsn) status

migrate-reset:
	goose -dir migrations postgres $(dsn) reset

migrate-create:
	goose -dir migrations create $(name) sql
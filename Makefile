.PHONY: swag

swag:
	swag init -g cmd/web/main.go --output cmd/web/docs


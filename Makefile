.PHONY: swag certs

swag:
	swag init -g cmd/web/main.go --output cmd/web/docs

certs:
	mkdir certs
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -pubout -out public.pem
	mv private.pem public.pem certs/
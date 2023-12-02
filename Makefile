#!make

swag:
	swag init -g http_server.go -d pkg/servers,pkg/core/response,internal/server/http/ctrl --parseDependency --parseInternal --output api

build-linux:
	swag init -g http_server.go -d pkg/servers,pkg/core/response,internal/server/http/ctrl --parseDependency --parseInternal --output api
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/app-amd64-linux cmd/app/http/main.go
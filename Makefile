#!make

swag:
	/Users/ndodanli/go/bin/swag init -g http_server.go -d pkg/servers,pkg/core/response,internal/server/http/controller --parseDependency --parseInternal --output api
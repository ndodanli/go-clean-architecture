package services

import jwtsvc "github.com/ndodanli/go-clean-architecture/pkg/infrastructure/jwt"

type AppServices struct {
	AuthService AuthServiceInterface
	JWTService  jwtsvc.JWTServiceInterface
}

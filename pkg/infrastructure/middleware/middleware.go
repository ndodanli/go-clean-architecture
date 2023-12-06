package mw

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/constant"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"strconv"
	"strings"
)

var (
	Authorize func(next echo.HandlerFunc) echo.HandlerFunc
	TraceID   func(next echo.HandlerFunc) echo.HandlerFunc
)

func Init(cfg *configs.Config) {
	Authorize = getJWTMiddleware(cfg, services.NewJWTService(cfg.Auth))
	TraceID = getTraceIDMiddleware()
}

func getJWTMiddleware(cfg *configs.Config, jwtService services.IJWTService) func(next echo.HandlerFunc) echo.HandlerFunc {
	validAudiences := strings.Split(cfg.Auth.JWT_AUDIENCES, ",")

	verifyAud := func(audiences []string) bool {
		if validAudiences[0] == "*" || audiences[0] == "*" {
			return true
		}
		for _, validAud := range validAudiences {
			for _, aud := range audiences {
				if validAud == aud {
					return true
				}
			}
		}
		return false
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := jwtService.ValidateToken(strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1))
			if err != nil {
				return httperr.UnauthorizedError
			}
			claims := token.Claims.(jwt.MapClaims)
			audiences, _ := claims.GetAudience()
			_ = audiences
			if !verifyAud(audiences) {
				return httperr.UnAuthorizedAudienceError
			}

			sub, _ := claims.GetSubject()
			subInt64, _ := strconv.ParseInt(sub, 10, 64)

			c.Set(constant.General.AuthUserId, &services.AuthUser{
				ID: subInt64,
			})

			return next(c)
		}
	}
}

func getTraceIDMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(constant.General.TraceIDKey, uuid.New().String())
			return next(c)
		}
	}
}

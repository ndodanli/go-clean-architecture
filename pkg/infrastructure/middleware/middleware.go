package mw

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/backend-api/configs"
	"github.com/ndodanli/backend-api/pkg/constant"
	httperr "github.com/ndodanli/backend-api/pkg/errors"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
	"strconv"
	"strings"
)

var (
	Auth    func(next echo.HandlerFunc) echo.HandlerFunc
	TraceID func(next echo.HandlerFunc) echo.HandlerFunc
)

func Init(cfg *configs.Config, appServices *services.AppServices, db *pgxpool.Pool) {
	Auth = getJWTMiddleware(cfg, appServices.JWTService, db)
	TraceID = getTraceIDMiddleware()
}

func getJWTMiddleware(cfg *configs.Config, jwtService services.IJWTService, db *pgxpool.Pool) func(next echo.HandlerFunc) echo.HandlerFunc {
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

			// Authorize
			requestedEndpoint := c.Path()
			var authorizeResponse *services.AuthorizeResponse
			authorizeResponse, err = jwtService.Authorize(c.Request().Context(), db, subInt64, requestedEndpoint, c.Request().Method)
			if err != nil {
				return httperr.ErrorWhileAuthorizingError
			}

			if !authorizeResponse.IsAuthorized {
				return httperr.UnauthorizedError
			}

			if authorizeResponse.IsBlocked {
				return httperr.UserBlockedError
			}

			c.Set(constant.General.AuthUserId, subInt64)
			c.Set(constant.General.AuthUserRoleIds, authorizeResponse.AppUserRoleIds)

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

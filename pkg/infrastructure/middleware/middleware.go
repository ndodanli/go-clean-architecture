package mw

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/constant"
	httperr "github.com/ndodanli/go-clean-architecture/pkg/errors"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
	"strconv"
	"strings"
)

var (
	Authenticate func(next echo.HandlerFunc) echo.HandlerFunc
	TraceID      func(next echo.HandlerFunc) echo.HandlerFunc
)

func Init(cfg *configs.Config, appServices *services.AppServices, db *pgxpool.Pool) {
	Authenticate = getJWTMiddleware(cfg, appServices.JWTService, db)
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
			requestedEndpoint := c.Path()
			_ = requestedEndpoint
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

			// Authorize
			exists := false
			// TODO: make this db function
			err = db.QueryRow(c.Request().Context(), `
WITH expanded_roles AS (SELECT UNNEST(roles) AS role_id
                        FROM app_user
                        WHERE id = $1
                          AND deleted_at = '0001-01-01T00:00:00Z'
                        LIMIT 1),
	-- Fetch endpoint ID directly
     endpoint AS (SELECT id
                     FROM endpoint
                     WHERE name = $2
                       AND deleted_at = '0001-01-01T00:00:00Z'
                     LIMIT 1)

	-- Check authorization
SELECT EXISTS (SELECT 1
               FROM expanded_roles er
                        JOIN role r ON r.id = er.role_id
               WHERE r.deleted_at = '0001-01-01T00:00:00Z'
                 AND (SELECT id from endpoint) = ANY (r.endpoint_ids));`, subInt64, requestedEndpoint).Scan(&exists)
			if err != nil {
				return httperr.UnauthorizedError
			}

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

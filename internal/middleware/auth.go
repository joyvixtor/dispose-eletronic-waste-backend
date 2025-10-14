package middleware

import (
	"net/http"
	"strings"

	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/adapters/auth"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			token := parts[1]
			claims, err := auth.ValidateToken(token, jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			c.Set("user_id", claims.UserId)
			c.Set("email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

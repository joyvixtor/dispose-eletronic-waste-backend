package api

import (
	auth "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/api/auth"
	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	Auth *auth.AuthHandler
}

func SetupRouter(e *echo.Echo, cfg RouteConfig) {
	auth := e.Group("/auth")
	auth.POST("/register", cfg.Auth.Register)
	auth.POST("/login", cfg.Auth.Login)
}

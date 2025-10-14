package api

import (
	"net/http"

	userauth "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/usecases/user-auth"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	useCase userauth.AuthUseCase
}

func NewAuthHandler(useCase userauth.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req userauth.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// if err := c.Validate(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"error": err.Error(),
	// 	})
	// }

	ctx := c.Request().Context()
	res, err := h.useCase.RegisterUser(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req userauth.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// if err := c.Validate(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"error": err.Error(),
	// 	})
	// }

	ctx := c.Request().Context()
	res, err := h.useCase.LoginUser(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

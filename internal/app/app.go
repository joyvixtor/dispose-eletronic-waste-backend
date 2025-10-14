package app

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/config"
	"github.com/labstack/echo/v4"

	user "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/adapters/repository"
	user_auth_handler "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/api/auth"
	user_auth "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/usecases/user-auth"

	dbpkg "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/adapters/database"
	api "github.com/joyvixtor/dispose-eletronic-waste-backend/internal/api"
)

type App struct {
	server *echo.Echo
	config *config.Config
	db     *sql.DB
	addr   string
}

func NewApp(cfg *config.Config) (*App, error) {
	e := echo.New()

	addr := fmt.Sprintf("0.0.0.0:%s", cfg.Port)

	db, err := dbpkg.New(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	return &App{
		server: e,
		config: cfg,
		db:     db.DB,
		addr:   addr,
	}, nil
}

func (a *App) Run() error {
	if err := a.initDependencies(); err != nil {
		return fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	slog.Info("Starting server", slog.String("addr", a.addr))
	return a.server.Start(a.addr)
}

func (a *App) initDependencies() error {
	//Repositories
	userRepo := user.NewUserRepository(a.db)

	//useCases
	userAuthUseCase := user_auth.NewAuthService(
		userRepo,
		a.config.JwtSecret,
		24,
	)

	//Handlers
	authHandler := user_auth_handler.NewAuthHandler(userAuthUseCase)

	api.SetupRouter(a.server, api.RouteConfig{
		Auth: authHandler,
	})

	return nil
}

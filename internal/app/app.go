package app

import (
	"context"
	"fmt"
	"time"

	"web/internal/auth"
	"web/internal/config"
	"web/internal/platform/database"
	platformlogger "web/internal/platform/logger"
	"web/internal/platform/openapi"
	"web/internal/role"
	"web/internal/user"
)

func Run(cfg config.Config) error {
	ctx := context.Background()
	logger := platformlogger.New(cfg.Environment)
	postgresDB, err := database.OpenPostgres(ctx, cfg.PostgresURL)
	if err != nil {
		return err
	}
	defer postgresDB.Close()

	roleStore := role.NewSQLStore(postgresDB)
	roleService := role.NewService(roleStore)
	roleHandler := role.NewHandler(roleService, logger)

	userStore := user.New(postgresDB)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	authRepository := auth.NewUserRepository(postgresDB)
	authService := auth.NewAuthService(authRepository, cfg.JWTSecret, 24*time.Hour)
	authHandler := auth.NewAuthHandler(authService)

	docsHandler := openapi.NewHandler()
	server := newServer(cfg.Address(), routes(docsHandler, roleHandler, userHandler, authHandler))
	logger.Info("server started", "address", server.Addr)
	return fmt.Errorf("serve HTTP: %w", server.ListenAndServe())
}

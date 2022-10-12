package main

import (
	"app-ez-pwd/internal/apis"
	"app-ez-pwd/internal/logger"
	"app-ez-pwd/internal/settings"
	"app-ez-pwd/internal/storage"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	settings.LoadConfiguration()
	defer logger.Logger.Sync()

	storage.ApplicationDB = storage.PrepareApplicationDB(settings.Settings.DatabaseURL)

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/api/v1/auth", apis.DoAuthPOST)
	e.DELETE("/api/v1/auth", apis.DoLogoutDELETE)
	e.POST("/api/v1/create-account", apis.CreateNewAccountPOST)

	apiV1 := e.Group("/api/v1")
	apiV1.Use(apis.VerifyAuthTokenMiddleware(""))
	apis.RouteUserSecretsApiHandlers(apiV1)

	go func() {
		signalStop := make(chan os.Signal, 1)
		signal.Notify(signalStop, syscall.SIGTERM, syscall.SIGINT)
		<-signalStop

		if err := e.Shutdown(context.Background()); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	if err := e.Start(settings.Settings.ApiHostPort); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	} else {
		logger.Logger.Info("api stopped")
	}
}

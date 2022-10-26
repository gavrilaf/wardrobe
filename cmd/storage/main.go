package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"

	"github.com/gavrilaf/wardrobe/pkg/utils/log"
	"github.com/gavrilaf/wardrobe/pkg/utils/server"
)

func main() {
	ctx := context.Background()
	log.InitLog(true)
	logger := log.FromContext(ctx)

	e := echo.New()
	e.Use(em.CORSWithConfig(em.DefaultCORSConfig))
	e.Use(em.Recover())

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	s := &http.Server{
		Addr:    ":8653",
		Handler: e,
	}

	logger.Infof("wardrobe is starting")

	quitChan := make(chan struct{}, 1)
	s = server.RunHTTP(ctx, s, quitChan)

	server.GracefulShutdown(ctx, quitChan, func(ctx context.Context) {
		server.ShutdownHTTP(ctx, s)
		logger.Info("wardrobe shutdown")
	})
}

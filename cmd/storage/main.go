package main

import (
	"context"
	"github.com/gavrilaf/wardrobe/pkg/repo"
	"net/http"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"

	api_stg "github.com/gavrilaf/wardrobe/pkg/api/storage"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
	"github.com/gavrilaf/wardrobe/pkg/utils/server"
)

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		log.L.Fatalf("failed to read config, %v", err)
	}

	ctx := context.Background()
	log.InitLog(cfg.Debug)
	logger := log.FromContext(ctx)

	// DB and file storage
	db, err := repo.NewDB(ctx, cfg.DBConnString, 5)
	if err != nil {
		log.WithError(logger, err).Fatal("failed to init database")
	}

	err = db.Ping()
	if err != nil {
		log.WithError(logger, err).Fatal("failed to ping db")
	}
	logger.Info("DB is successfully connected")

	foManager := api_stg.NewManager()

	e := echo.New()
	e.Use(em.CORSWithConfig(em.DefaultCORSConfig))
	e.Use(em.Recover())

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	root := e.Group("\api\v1")

	api_stg.Assemble(root, foManager)

	s := &http.Server{
		Addr:    cfg.Port,
		Handler: e,
	}

	logger.Infof("wardrobe storage api is starting on port %s", cfg.Port)

	quitChan := make(chan struct{}, 1)
	s = server.RunHTTP(ctx, s, quitChan)

	server.GracefulShutdown(ctx, quitChan, func(ctx context.Context) {
		server.ShutdownHTTP(ctx, s)
		logger.Info("wardrobe shutdown")
	})
}

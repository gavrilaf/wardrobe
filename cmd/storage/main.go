package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"

	mw "github.com/gavrilaf/wardrobe/pkg/api/middleware"
	api_stg "github.com/gavrilaf/wardrobe/pkg/api/storage"
	"github.com/gavrilaf/wardrobe/pkg/fs/minio"
	"github.com/gavrilaf/wardrobe/pkg/repo"
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

	// Database
	db, err := repo.NewDB(ctx, cfg.DBConnString, 5)
	if err != nil {
		log.WithError(logger, err).Fatal("failed to init database")
	}

	err = db.Ping()
	if err != nil {
		log.WithError(logger, err).Fatal("failed to ping db")
	}
	logger.Info("DB is successfully connected")

	if err = db.Migrate(ctx, "./migration"); err != nil {
		log.WithError(logger, err).Fatal("DB migration failed")
	}

	logger.Info("DB migration ok")

	// Storage
	stg, err := minio.New(cfg.MinioEndpoint, cfg.MinioUser, cfg.MinioPassword, cfg.FOBucket)
	if err != nil {
		log.WithError(logger, err).Fatal("failed to connect to the files storage")
	}

	if err = stg.Prepare(ctx); err != nil {
		log.WithError(logger, err).Fatal("failed to connect to prepare storage")
	}

	// API

	foManager := api_stg.NewManager(api_stg.Config{
		Tx:          db,
		FileObjects: db,
		Tags:        db,
		Stg:         stg,
	})

	e := echo.New()
	e.Use(em.CORSWithConfig(em.DefaultCORSConfig))
	e.Use(em.Recover())
	e.Use(mw.Measure)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	root := e.Group("/api/v1")

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

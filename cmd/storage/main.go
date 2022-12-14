package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	em "github.com/labstack/echo/middleware"

	mw "github.com/gavrilaf/wardrobe/pkg/api/middleware"
	apistg "github.com/gavrilaf/wardrobe/pkg/api/storage"
	"github.com/gavrilaf/wardrobe/pkg/domain/stglogic"
	"github.com/gavrilaf/wardrobe/pkg/fs/minio"
	"github.com/gavrilaf/wardrobe/pkg/repo"
	"github.com/gavrilaf/wardrobe/pkg/utils/idgen"
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

	// Storage
	stg, err := minio.New(cfg.MinioEndpoint, cfg.MinioUser, cfg.MinioPassword)
	if err != nil {
		log.WithError(logger, err).Fatal("failed to create storage")
	}

	if err = stg.Ping(); err != nil {
		log.WithError(logger, err).Fatal("storage is offline")
	}

	logger.Info("storage is online")

	// Configurator
	nodeID, err := idgen.NodeID()
	if err != nil {
		log.WithError(logger, err).Fatal("failed to retrieve node id")
	}

	snowflake := idgen.NewSnowflake(nodeID)

	stgConfigurator := stglogic.NewConfigurator(stg, snowflake)

	err = stgConfigurator.PrepareStorage(ctx)
	if err != nil {
		log.WithError(logger, err).Fatal("failed to prepare storage")
	}

	logger.Infof("storage is ready, configurator created with node id %d", nodeID)

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

	// API

	foManager := apistg.NewManager(apistg.Config{
		Tx:              db,
		InfoObjects:     db,
		FS:              stg,
		StgConfigurator: stgConfigurator,
	})

	e := echo.New()
	e.Use(em.CORSWithConfig(em.DefaultCORSConfig))
	e.Use(em.Recover())
	e.Use(mw.Measure)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	root := e.Group("/api/v1")

	apistg.Assemble(root, foManager)

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

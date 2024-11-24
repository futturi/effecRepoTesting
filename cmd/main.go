package main

import (
	"context"
	"effectiveMobile/internal/config"
	"effectiveMobile/internal/handler"
	"effectiveMobile/internal/logger"
	"effectiveMobile/internal/server"
	"effectiveMobile/internal/service"
	"effectiveMobile/internal/store"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// @title Song API
// @description API Server 4 Song

// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.NewConfig()
	log := logger.InitLogger(cfg.LogLevel)
	ctx := logger.ContextWithLogger(context.Background(), log)
	if err != nil {
		log.Errorw("error with creating config", zap.Error(err))
		return
	}
	db, err := store.InitDB(ctx, cfg.DbConnectionString)
	if err != nil {
		log.Errorw("error with creating db", zap.Error(err))
		return
	}
	storeLevel := store.NewStore(db)
	log.Debug("created store level")
	service := service.NewService(&storeLevel, cfg.ApiUrl, cfg.ApiTimeout)
	log.Debug("created service level")
	handl := handler.NewHandler(service)
	log.Debug("created handler level")
	r := handl.InitRoutes(ctx)
	go func() {
		log.Debugw("server is running", "port", cfg.Port)
		if err := server.NewServer(r, cfg.Port); err != nil {
			log.Errorw("error with starting server", zap.Error(err))
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err := store.ShutdownDb(ctx, db); err != nil {
		log.Errorw("error with shutting down service", zap.Error(err))
		os.Exit(1)
	}

	log.Debugw("server is shutting down")
}

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"
	"url-shortener/app/config"
	"url-shortener/app/handlers"
	"url-shortener/app/repository"
	"url-shortener/cmd/server"
	"url-shortener/pkg/database"
	"url-shortener/pkg/logger"
)

func main() {
	log := logger.SetupLogger()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Error("error init config:", logger.Err(err))
	}

	db, err := database.NewDB(&cfg.DBConfig)
	if err != nil {
		log.Error("error init database:", logger.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	storage := repository.NewURLRepository(db)
	h := handlers.InitRoutes(log, storage)

	srv := &server.Server{}

	if err := srv.Run(&cfg.HTTPServer, h); err != nil {
		log.Error("Cant run Server", logger.Err(err))
		os.Exit(1)
	}

	log.Info("Server Started on", slog.String("port", cfg.HTTPServer.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info("Shutting down server")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.ShutDown(shutDownCtx); err != nil {
		log.Error("Cant shutdown server", logger.Err(err))
	}

}

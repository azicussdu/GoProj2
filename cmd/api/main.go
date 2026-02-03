package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/azicussdu/GoProj2/internal/config"
	"github.com/azicussdu/GoProj2/internal/handler"
	"github.com/azicussdu/GoProj2/internal/pkg/logger"
	"github.com/azicussdu/GoProj2/internal/repository"
	"github.com/azicussdu/GoProj2/internal/server"
	"github.com/azicussdu/GoProj2/internal/service"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	slogger := logger.New(cfg.LogLevel)
	slog.SetDefault(slogger)

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		return
	}

	courseRepo := repository.NewPsqCourseRepo(db)
	courseService := service.NewCourseService(courseRepo)

	h := handler.NewHandler(courseService)
	router, err := h.InitRoutes()

	srv := server.New(router, cfg.Port)
	err = srv.Run()
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

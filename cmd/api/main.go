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
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	slogger := logger.New(cfg.LogLevel)
	slog.SetDefault(slogger)

	router, err := buildApp(cfg)
	if err != nil {
		slog.Error("failed to build app", "error", err.Error())
		os.Exit(1)
	}

	srv := server.New(router, cfg.Port)
	err = srv.Run()
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("Server started", "port", cfg.Port)
}

func buildApp(cfg *config.Config) (*gin.Engine, error) {
	// создаем db - с его помощью будем делать запросы в БД Postgres
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("Error with DB connection")
		return nil, err
	}

	// Тут только создаем репозитории
	courseRepo := repository.NewPsqCourseRepo(db)
	lessonRepo := repository.NewPsgLessonRepo(db)
	userRepo := repository.NewPsgUserRepo(db)

	// Cобрали все сервисе в одним файле
	services := &service.Services{
		Course: service.NewCourseService(courseRepo, lessonRepo),
		Lesson: service.NewLessonService(lessonRepo, courseRepo),
		Auth:   service.NewAuthService(userRepo),
	}

	h := handler.NewHandler(services)
	router, err := h.InitRoutes()

	return router, nil
}

package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/azicussdu/GoProj2/internal/auth"
	"github.com/azicussdu/GoProj2/internal/config"
	"github.com/azicussdu/GoProj2/internal/handler"
	"github.com/azicussdu/GoProj2/internal/pkg/logger"
	"github.com/azicussdu/GoProj2/internal/repository"
	"github.com/azicussdu/GoProj2/internal/server"
	"github.com/azicussdu/GoProj2/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	enrollmentRepo := repository.NewPsgEnrollmentRepo(db)
	userRepo := repository.NewPsgUserRepo(db)
	redisClient := initRedisClient(cfg)

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, cfg.JWT.Issuer)

	// Cобрали все сервисе в одним файле
	services := &service.Services{
		Course:     service.NewCourseService(courseRepo, lessonRepo, enrollmentRepo, db, redisClient),
		Lesson:     service.NewLessonService(lessonRepo, courseRepo, db),
		Enrollment: service.NewEnrollmentService(enrollmentRepo, courseRepo),
		Auth:       service.NewAuthService(userRepo, jwtManager),
	}

	h := handler.NewHandler(services, jwtManager)
	router, err := h.InitRoutes()

	return router, nil
}

func initRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("redis is unavailable, continuing without cache", "error", err.Error())
		_ = client.Close()
		return nil
	}

	slog.Info("Redis connected successfully")
	return client
}

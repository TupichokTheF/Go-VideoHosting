package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"project/internal/application/services"
	"project/internal/core"
	app_kafka "project/internal/infrastructure/brokers/kafka"
	"project/internal/infrastructure/cache"
	"project/internal/infrastructure/database/postgres"
	app_redis "project/internal/infrastructure/database/redis"
	infra_messaging "project/internal/infrastructure/messaging"
	"project/internal/infrastructure/repositories"
	"project/internal/infrastructure/security"
	"project/internal/infrastructure/storage"
	http_handlers "project/internal/presentation/http/handlers"
	"project/internal/presentation/http/routers"
	"syscall"
	"time"
)

// @title    Go Video Hosting API
// @version  1.0
// @host     localhost:8080
// @BasePath /api/v1
func main() {
	start()
}

func start() {
	cfg := core.LoadConfig()

	db, err := postgres.OpenConnection(cfg.DataBaseConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection()

	redisClient, err := app_redis.NewClient(&cfg.RedisConfig)
	if err != nil {
		log.Fatal(err)
	}

	minioClient, err := storage.NewMinioClient(&cfg.MinioConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger := core.SetupLogger()
	slog.SetDefault(logger)

	userRepo := repositories.NewUserRepository(db.ConnPool)
	videoRepo := repositories.NewVideoRepository(db.ConnPool)

	tokenCache := cache.NewTokenCache(redisClient, cfg.RefreshTTL)

	jwtManager := security.NewJWTManager(cfg.AccessSecretKey, cfg.RefreshSecretKey, cfg.AccessTTL, cfg.RefreshTTL)
	hasher := security.NewBcryptHasher()
	videoStorage := storage.NewMinioService(minioClient, cfg.MinioConfig.Bucket, cfg.MinioConfig.TTL)

	kafkaProducer := app_kafka.NewProducer(&cfg.KafkaConfig)
	defer func() {
		err := app_kafka.CloseProducer(kafkaProducer)
		logger.Error("closing of kafka producer completed with error", "error", err)
	}()
	kafkaPublisher := infra_messaging.NewPublisher(kafkaProducer)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, jwtManager, hasher, tokenCache)
	videoService := services.NewVideoService(videoRepo, videoStorage, kafkaPublisher)

	authHandler := http_handlers.NewAuthHandler(authService)
	userHandler := http_handlers.NewUserHandler(userService)
	videoHandler := http_handlers.NewVideoHandler(videoService)

	routersOptions := []routers.Option{
		routers.WithAuthRouter(authHandler),
		routers.WithUserRouter(userHandler, authService),
		routers.WithVideoRouter(videoHandler, authService),
	}
	if cfg.Swagger {
		routersOptions = append(routersOptions, routers.WithSwagger())
	}

	router := routers.GetRouter(routersOptions...)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	server := http.Server{Addr: cfg.GetAddress(), Handler: router}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("error while listening server", "error", err)
			stop()
		}
	}()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("error while shuting down server", "error", err)
	}
}

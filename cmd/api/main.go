package main

import (
	"log"
	"net/http"
	"project/internal/application/services"
	"project/internal/core"
	"project/internal/infrastructure/cache"
	"project/internal/infrastructure/database/postgres"
	app_redis "project/internal/infrastructure/database/redis"
	"project/internal/infrastructure/repositories"
	"project/internal/infrastructure/security"
	"project/internal/infrastructure/storage"
	"project/internal/presentation/handlers"
	"project/internal/presentation/routers"
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

	userRepo := repositories.NewUserRepository(db.ConnPool)
	videoRepo := repositories.NewVideoRepository(db.ConnPool)

	tokenCache := cache.NewTokenCache(redisClient, cfg.RefreshTTL)

	jwtManager := security.NewJWTManager(cfg.AccessSecretKey, cfg.RefreshSecretKey, cfg.AccessTTL, cfg.RefreshTTL)
	hasher := security.NewBcryptHasher()
	videoStorage := storage.NewMinioService(minioClient, cfg.MinioConfig.Bucket, cfg.MinioConfig.TTL)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, jwtManager, hasher, tokenCache)
	videoService := services.NewVideoService(videoRepo, videoStorage)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	videoHandler := handlers.NewVideoHandler(videoService)

	routersOptions := []routers.Option{
		routers.WithAuthRouter(authHandler),
		routers.WithUserRouter(userHandler, authService),
		routers.WithVideoRouter(videoHandler, authService),
	}

	if cfg.Swagger {
		routersOptions = append(routersOptions, routers.WithSwagger())
	}

	router := routers.GetRouter(routersOptions...)
	if err := http.ListenAndServe(cfg.GetAddress(), router); err != nil {
		log.Fatal(err.Error())
	}
}

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
	"project/internal/presentation/handlers"
	"project/internal/presentation/routers"
)

// @title    Go Messenger API
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

	userRepo := repositories.NewUserRepository(db.ConnPool)

	tokenCache := cache.NewTokenCache(redisClient, cfg.RefreshTTL)

	jwtManager := security.NewJWTManager(cfg.AccessSecretKey, cfg.RefreshSecretKey, cfg.AccessTTL, cfg.RefreshTTL)
	hasher := security.NewBcryptHasher()

	userService := services.NewUserService(userRepo, tokenCache)
	authService := services.NewAuthService(userRepo, jwtManager, hasher, tokenCache)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	routersOptions := []routers.Option{
		routers.WithAuthRouter(authHandler),
		routers.WithUserRouter(userHandler, authService),
	}

	if cfg.Swagger {
		routersOptions = append(routersOptions, routers.WithSwagger())
	}

	router := routers.GetRouter(routersOptions...)
	if err := http.ListenAndServe(cfg.GetAddress(), router); err != nil {
		log.Fatal(err.Error())
	}
}

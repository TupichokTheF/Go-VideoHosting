package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

type dependencies struct {
	db            *postgres.Connection
	redisClient   *redis.Client
	minioClient   *minio.Client
	kafkaProducer *kafka.Writer
}

type appServices struct {
	userService  *services.UserService
	authService  *services.AuthService
	videoService *services.VideoService
}

type Application struct {
	*dependencies
	*appServices
	onClose []io.Closer
	cfg     *core.Config
}

func (app *Application) Close() error {
	var errs []error
	for _, closer := range app.onClose {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func SetupApplication(cfg *core.Config) (*Application, error) {
	app := &Application{cfg: cfg, onClose: make([]io.Closer, 0)}

	if err := app.setupDeps(); err != nil {
		return nil, errors.Join(err, app.Close())
	}

	app.setupServices()

	return app, nil
}

func (app *Application) setupDeps() error {
	db, err := postgres.OpenConnection(app.cfg.DataBaseConfig)
	if err != nil {
		return fmt.Errorf("error while setting up postgres: %w", err)
	}
	app.onClose = append(app.onClose, db)

	redisClient, err := app_redis.NewClient(&app.cfg.RedisConfig)
	if err != nil {
		return fmt.Errorf("error while setting up redis: %w", err)
	}
	app.onClose = append(app.onClose, redisClient)

	minioClient, err := storage.NewMinioClient(&app.cfg.MinioConfig)
	if err != nil {
		return fmt.Errorf("error while setting up minio: %w", err)
	}

	kafkaProducer := app_kafka.NewProducer(&app.cfg.KafkaConfig)
	app.onClose = append(app.onClose, kafkaProducer)

	app.dependencies = &dependencies{
		db:            db,
		redisClient:   redisClient,
		minioClient:   minioClient,
		kafkaProducer: kafkaProducer,
	}

	return nil
}

func (app *Application) setupServices() {
	userRepo := repositories.NewUserRepository(app.db.ConnPool)
	videoRepo := repositories.NewVideoRepository(app.db.ConnPool)

	tokenCache := cache.NewTokenCache(app.redisClient, app.cfg.RefreshTTL)

	jwtManager := security.NewJWTManager(app.cfg.AccessSecretKey, app.cfg.RefreshSecretKey, app.cfg.AccessTTL, app.cfg.RefreshTTL)
	hasher := security.NewBcryptHasher()
	videoStorage := storage.NewMinioService(app.minioClient, app.cfg.MinioConfig.Bucket, app.cfg.MinioConfig.TTL)

	kafkaPublisher := infra_messaging.NewPublisher(app.kafkaProducer)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, jwtManager, hasher, tokenCache)
	videoService := services.NewVideoService(videoRepo, videoStorage, kafkaPublisher)

	app.appServices = &appServices{
		userService:  userService,
		authService:  authService,
		videoService: videoService,
	}
}

func (app *Application) Run(ctx context.Context) error {
	gErr, ctx := errgroup.WithContext(ctx)
	gErr.Go(func() error {
		return app.startHTTPServer(ctx)
	})

	return gErr.Wait()
}

func (app *Application) startHTTPServer(ctx context.Context) error {
	router := app.setupRouter()
	server := http.Server{Addr: app.cfg.GetAddress(), Handler: router}

	errStream := make(chan error, 1)
	go func() {
		errStream <- server.ListenAndServe()
	}()

	select {
	case err := <-errStream:
		return fmt.Errorf("error while listening the server: %w", err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("error while shutting down the server: %w", err)
	}

	return nil
}

func (app *Application) setupRouter() *chi.Mux {
	authHandler := http_handlers.NewAuthHandler(app.authService)
	userHandler := http_handlers.NewUserHandler(app.userService)
	videoHandler := http_handlers.NewVideoHandler(app.videoService)

	routersOptions := []routers.Option{
		routers.WithAuthRouter(authHandler),
		routers.WithUserRouter(userHandler, app.authService),
		routers.WithVideoRouter(videoHandler, app.authService),
	}
	if app.cfg.Swagger {
		routersOptions = append(routersOptions, routers.WithSwagger())
	}

	return routers.GetRouter(routersOptions...)
}

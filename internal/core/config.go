package core

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	BasePath string `env:"BASE_PATH" env-required:"true"`
	HTTPConfig
	DataBaseConfig
	RedisConfig
	LocalConfig
	JWTConfig
	MinioConfig
}

type HTTPConfig struct {
	HTTPHost string `env:"HTTP_HOST"`
	HTTPPort int    `env:"HTTP_PORT"`
}

func (conf *HTTPConfig) GetAddress() string {
	return fmt.Sprintf("%s:%v", conf.HTTPHost, conf.HTTPPort)
}

type DataBaseConfig struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     int    `env:"DB_PORT" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
}

func (d *DataBaseConfig) GetURL() string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%v/%s",
		d.User, d.Password, d.Host, d.Port, d.Name)

	return url
}

type RedisConfig struct {
	Host string `env:"REDIS_HOST" env-required:"true"`
	Port int    `env:"REDIS_PORT" env-required:"true"`
}

func (cfg *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
}

type LocalConfig struct {
	Swagger bool `env:"SWAGGER" env-required:"false"`
}

type JWTConfig struct {
	AccessSecretKey  []byte        `env:"ACCESS_SECRET_KEY" env-required:"true"`
	RefreshSecretKey []byte        `env:"REFRESH_SECRET_KEY" env-required:"true"`
	AccessTTL        time.Duration `env:"ACCESS_TOKEN_TTL" env-required:"true"`
	RefreshTTL       time.Duration `env:"REFRESH_TOKEN_TTL" env-required:"true"`
}

type MinioConfig struct {
	Host     string `env:"MINIO_HOST" env-required:"true"`
	Port     int    `env:"MINIO_PORt" env-required:"true"`
	Bucket   string `env:"MINIO_BUCKET" env-required:"true"`
	User     string `env:"MINIO_USER" env-required:"true"`
	Password string `env:"MINIO_PASSWORD" env-required:"true"`
	UseSSL   bool   `env:"MINIO_SSL" env-required:"true"`
}

func (cfg *MinioConfig) Endpoint() string {
	return fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Can't load env file")
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Can't read config file")
	}

	return &cfg
}

package core

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	KafkaConfig
	TranscoderConfig
	MinioConfig
}

type KafkaConfig struct {
	Host    string   `env:"KAFKA_HOST" env-default:"localhost"`
	Port    int      `env:"KAFKA_PORT" env-default:"9092"`
	GroupID string   `env:"KAFKA_GROUP_ID" env-default:"transcoding"`
	Topics  []string `env:"KAFKA_TOPICS" env-separator:"," env-default:"video.events"`
}

func (cfg *KafkaConfig) Address() string {
	return fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
}

type TranscoderConfig struct {
	Bin string `env:"TRANSCODER_BIN" env-default:"/usr/bin/ffmpeg"`
}

type MinioConfig struct {
	Host     string        `env:"MINIO_HOST" env-default:"localhost"`
	Port     int           `env:"MINIO_PORT" env-default:"9000"`
	Bucket   string        `env:"MINIO_BUCKET" env-default:"videos"`
	User     string        `env:"MINIO_USER" env-required:"true"`
	Password string        `env:"MINIO_PASSWORD" env-required:"true"`
	TTL      time.Duration `env:"MINIO_TTL" env-default:"30m"`
	UseSSL   bool          `env:"MINIO_SSL" env-default:"false"`
}

func (cfg *MinioConfig) Endpoint() string {
	return fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("no .env file, reading from environment")
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal("Can't read config file")
	}

	return &cfg
}

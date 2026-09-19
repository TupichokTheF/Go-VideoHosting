package core

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	KafkaConfig
	TranscoderConfig
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

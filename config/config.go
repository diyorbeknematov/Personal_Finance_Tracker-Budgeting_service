package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	GRPC_PORT      int    `yaml:"grpc_port"`
	Redis_HOST     string `yaml:"redis_host"`
	Redis_PORT     int    `yaml:"redis_port"`
	Redis_PASSWORD string `yaml:"redis_password"`
	Redis_DB       int    `yaml:"redis_db"`
	MONGODB_NAME   string `yaml:"mongodb_name"`
	MONGODB_URI    string `yaml:"mongodb_uri"`
}

func Load() *Config {
	if err := godotenv.Load("./../../.env"); err != nil {
		log.Println("Error loading .env file")
	}

	config := &Config{}

	config.GRPC_PORT = cast.ToInt(coalesce("GRPC_PORT", 50051))

	config.Redis_HOST = cast.ToString(coalesce("REDIS_HOST", "localhost"))
	config.Redis_PORT = cast.ToInt(coalesce("REDIS_PORT", 6379))
	config.Redis_PASSWORD = cast.ToString(coalesce("REDIS_PASSWORD", ""))
	config.Redis_DB = cast.ToInt(coalesce("REDIS_DB", 0))

	config.MONGODB_NAME = cast.ToString(coalesce("MONGODB_NAME", "mongo"))
	config.MONGODB_URI = cast.ToString(coalesce("MONGODB_URI", "mongodb://mongo:27017"))

	return config
}

func coalesce(key string, defaults interface{}) interface{} {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaults
}

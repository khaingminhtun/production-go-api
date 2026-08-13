package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string

	ServerPort string

	Loglevel string

	DB DatabaseConfig

	Redis RedisConfig

	SendGrid SendGridConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type SendGridConfig struct {
	APIKey string

	FromEmail string

	FromName string
}

func Load() *Config {

	err := godotenv.Load("configs/development.env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5433")
	viper.SetDefault("DB_USER", "production")
	viper.SetDefault("DB_NAME", "production_api")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)

	return &Config{

		AppEnv: viper.GetString("APP_ENV"),

		ServerPort: viper.GetString("SERVER_PORT"),

		Loglevel: viper.GetString("LOG_LEVEL"),

		DB: DatabaseConfig{

			Host: viper.GetString("DB_HOST"),

			Port: viper.GetString("DB_PORT"),

			User: viper.GetString("DB_USER"),

			Password: viper.GetString("DB_PASSWORD"),

			Name: viper.GetString("DB_NAME"),

			SSLMode: viper.GetString("DB_SSL_MODE"),
		},

		Redis: RedisConfig{
			Host: viper.GetString("REDIS_HOST"),

			Port: viper.GetString("REDIS_PORT"),

			Password: viper.GetString("REDIS_PASSWORD"),

			DB: viper.GetInt("REDIS_DB"),
		},

		SendGrid: SendGridConfig{
			APIKey: viper.GetString("SENDGRID_API_KEY"),

			FromEmail: viper.GetString("SENDGRID_FROM_EMAIL"),

			FromName: viper.GetString("SENDGRID_FROM_NAME"),
		},
	}

}

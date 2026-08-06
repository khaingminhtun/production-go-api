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
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {

	err := godotenv.Load("configs/development.env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	viper.AutomaticEnv()

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
	}

}

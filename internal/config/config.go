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
}

func Load() *Config {

	err := godotenv.Load("configs/development.env")

	if err != nil{
		log.Println("Error loading .env file")
	}

	viper.AutomaticEnv()

	return &Config{

		AppEnv: viper.GetString("APP_ENV"),

		ServerPort: viper.GetString("SERVER_PORT"),

		Loglevel: viper.GetString("LOGLEVEL"),
	}

}
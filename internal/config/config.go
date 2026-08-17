package config

import (
	"log"
	"time"

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

	JwtEnv JWTConfig
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

type JWTConfig struct {
	Secret            string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
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

	viper.SetDefault("JWT_SECRET", "default_secret_please_change_this")
	viper.SetDefault("JWT_ACCESS_EXPIRATION", "15m")
	viper.SetDefault("JWT_REFRESH_EXPIRATION", "720h")

	accessExpiration, err := time.ParseDuration(viper.GetString("JWT_ACCESS_EXPIRATION"))
	if err != nil {
		log.Printf("Invalid JWT_ACCESS_EXPIRATION format. Falling back to 15m: %v", err)
		accessExpiration = 15 * time.Minute
	}

	refreshExpiration, err := time.ParseDuration(viper.GetString("JWT_REFRESH_EXPIRATION"))
	if err != nil {
		log.Printf("Invalid JWT_REFRESH_EXPIRATION format. Falling back to 720h: %v", err)
		refreshExpiration = 720 * time.Hour
	}

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

		JwtEnv: JWTConfig{
			Secret:            viper.GetString("JWT_SECRET"),
			AccessExpiration:  accessExpiration,  // Injected time.Duration
			RefreshExpiration: refreshExpiration, // Injected time.Duration
		},
	}

}

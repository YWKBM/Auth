package config

import (
	"os"
)

type Config struct {
	HOST string
	PORT string

	SECRET_KEY string

	DB_CONFIG PostgresConfig
	RABBIT_MQ RabbitMQConfig
	ADMIN     AdminConfig
}

type PostgresConfig struct {
	DB_USER  string
	DB_NAME  string
	SSL_MODE string
	DB_PORT  string
	DB_PASS  string
	DB_HOST  string
}

type RabbitMQConfig struct {
	RABBIT_URL string
}

type AdminConfig struct {
	LOGIN    string
	EMAIL    string
	PASSWORD string
}

func Init() *Config {
	host, _ := os.LookupEnv("HOST")
	port, _ := os.LookupEnv("PORT")
	secretkey, _ := os.LookupEnv("SECRET_KEY")

	dbUser, _ := os.LookupEnv("DB_USER")
	dbName, _ := os.LookupEnv("DB_NAME")
	sslMode, _ := os.LookupEnv("SSL_MODE")
	dbPort, _ := os.LookupEnv("DB_PORT")
	dbPass, _ := os.LookupEnv("DB_PASS")
	dbHost, _ := os.LookupEnv("DB_HOST")

	rabbitUrl, _ := os.LookupEnv("RABBIT_URL")

	admLogin, _ := os.LookupEnv("LOGIN")
	admEmail, _ := os.LookupEnv("EMAIL")
	admPass, _ := os.LookupEnv("PASSWORD")

	return &Config{
		HOST:       host,
		PORT:       port,
		SECRET_KEY: secretkey,

		DB_CONFIG: PostgresConfig{
			DB_USER:  dbUser,
			DB_NAME:  dbName,
			SSL_MODE: sslMode,
			DB_PORT:  dbPort,
			DB_PASS:  dbPass,
			DB_HOST:  dbHost,
		},

		RABBIT_MQ: RabbitMQConfig{
			RABBIT_URL: rabbitUrl,
		},

		ADMIN: AdminConfig{
			LOGIN:    admLogin,
			EMAIL:    admEmail,
			PASSWORD: admPass,
		},
	}
}

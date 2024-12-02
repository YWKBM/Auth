package main

import (
	"auth/config"
	"auth/database"
	"auth/handler"
	"auth/repo"
	"auth/services"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

// @title           Swagger Auth API
// @version         1.0
// @description     Auth API

// @host      localhost:8080
// @BasePath  /api/auth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Bearer-Token

func main() {
	config := config.Init()

	connectionInfo := database.ConnectionInfo{
		Username: config.DB_CONFIG.DB_USER,
		Host:     config.DB_CONFIG.DB_HOST,
		Port:     config.DB_CONFIG.DB_PORT,
		DBName:   config.DB_CONFIG.DB_NAME,
		SSLMode:  config.DB_CONFIG.SSL_MODE,
		Password: config.DB_CONFIG.DB_PASS,
	}

	db, err := database.NewPostgresConnection(connectionInfo)
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	defer db.Close()

	logger := logrus.New()

	repo := repo.NewRepos(db)
	servs := services.NewServices(repo, *config)
	handler := handler.NewHandler(servs, config, logger)

	logger.Info("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	//allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	if err == http.ListenAndServe(fmt.Sprintf("%s:%s", config.HOST, config.PORT), handler.Init()
	handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		// handlers.AllowCredentials(),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"}),
		handlers.AllowedMethods([]string{"OPTIONS", "POST", "GET", "DELETE", "PUT"}),
		// handlers.AllowCredentials(),
	)(handler.Init())); err != nil {
	if err != nil {
		log.Fatal()
	}

}

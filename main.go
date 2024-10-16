package main

import (
	"auth/database"
	"auth/handler"
	"auth/repo"
	"auth/services"
	"log"
	"net/http"
	"time"
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
	connectionInfo := database.ConnectionInfo{
		Username: "postgres",
		Host:     "localhost",
		Port:     "5432",
		DBName:   "Auth",
		SSLMode:  "disable",
		Password: "123456",
	}

	db, err := database.NewPostgresConnection(connectionInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repo := repo.NewRepos(db)
	servs := services.NewServices(repo)
	handler := handler.NewHandler(servs)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.Init(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}
}

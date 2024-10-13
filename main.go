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

	//handler.Init()

	// 	Host     string
	// Port     string
	// Username string
	// DBName   string
	// SSLMode  string
	// Password string
}

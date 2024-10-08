package main

import (
	"auth/database"
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

	for {
		time.Sleep(time.Second * 5)
		db.Ping()
	}
	// 	Host     string
	// Port     string
	// Username string
	// DBName   string
	// SSLMode  string
	// Password string
}

package handler

import (
	"auth/config"
	"auth/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	config      *config.Config
	log         *logrus.Logger
	authHandler *AuthHandler
}

func NewHandler(servs *services.Services, config *config.Config, logger *logrus.Logger) *Handler {
	return &Handler{
		config:      config,
		log:         logger,
		authHandler: newAuthHandler(servs.AuthService)}
}

func (h *Handler) Init() *mux.Router {
	router := mux.NewRouter()

	//router.Use(h.cors)
	router.Use(h.requestLogging)
	//router.Use(h.cors)

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowedMethods:   []string{"*"},
	// 	AllowCredentials: true,
	// })

	// router.Use(c.Handler)

	router.HandleFunc("/api/auth/sign_up", h.errorProcessing(h.authHandler.SignUp)).Methods("POST")
	router.HandleFunc("/api/auth/sign_in", h.errorProcessing(h.authHandler.SignIn)).Methods("POST")
	router.HandleFunc("/api/auth/renew", h.errorProcessing(h.authHandler.RenewCredentials)).Methods("POST")

	router.Handle("/api/auth/sign_out", h.userIdentity(http.HandlerFunc(h.errorProcessing(h.authHandler.SignOut)))).Methods("POST")
	router.Handle("/api/auth/change_password", h.userIdentity(http.HandlerFunc(h.errorProcessing(h.authHandler.ChangePassword)))).Methods("POST")

	return router
}

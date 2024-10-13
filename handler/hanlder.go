package handler

import (
	"auth/services"

	"github.com/gorilla/mux"
)

type Handler struct {
	authHandler *AuthHandler
}

func NewHandler(servs *services.Services) *Handler {
	return &Handler{authHandler: newAuthHandler(servs.AuthService)}
}

func (h *Handler) Init() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/auth/sign_up", h.authHandler.SignUp)
	router.HandleFunc("/api/auth/sign_in", h.authHandler.SignIn)
	router.HandleFunc("/api/auth/sign_out", h.authHandler.SignOut)
	router.HandleFunc("/api/auth/renew", h.authHandler.RenewCredentials)

	//router.Use(h.userIdentity)
	router.HandleFunc("/api/auth/change_password", h.authHandler.ChangePassword)

	return router
}

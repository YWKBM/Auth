package handler

import (
	"auth/services"
	"net/http"

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

	router.HandleFunc("/api/auth/sign_up", h.authHandler.SignUp).Methods("POST")
	router.HandleFunc("/api/auth/sign_in", h.authHandler.SignIn).Methods("POST")
	router.HandleFunc("/api/auth/renew", h.authHandler.RenewCredentials).Methods("POST")

	router.Handle("/api/auth/sign_out", h.userIdentity(http.HandlerFunc(h.authHandler.SignOut))).Methods("POST")
	router.Handle("/api/auth/change_password", h.userIdentity(http.HandlerFunc(h.authHandler.ChangePassword))).Methods("POST")

	return router
}

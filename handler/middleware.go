package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	authHeader = "Authorization"
	userCtx    = "UserId"
)

func (h *Handler) userIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get(authHeader)
		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := h.authHandler.authService.ParseAccessToken(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Add(userCtx, fmt.Sprintf("%v", userId))
		next.ServeHTTP(w, r)
	})

}

func getUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	id := r.Header.Get("UserId")

	if id == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return 0, errors.New("Unauthorized")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 0, err
	}

	return idInt, nil
}

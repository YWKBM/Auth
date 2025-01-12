package handler

import (
	"auth/customErrors"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	authHeader = "Authorization"
	userCtx    = "UserId"
)

func (h *Handler) requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, _ := io.ReadAll(r.Body)

		h.log.Info(fmt.Printf("Request body: %v", string(buf)))

		reader := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = reader

		next.ServeHTTP(w, r)
	})
}

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h *Handler) errorProcessing(f ErrorHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			switch err.(type) {
			case *customErrors.ValidationError:
				h.log.Error(err)
				w.WriteHeader(400)
				w.Write([]byte(err.Error()))
			case *customErrors.NotFoundError:
				h.log.Error(err)
				w.WriteHeader(404)
				w.Write([]byte(err.Error()))
			case *customErrors.AlreadyExistsError:
				h.log.Error(err)
				w.WriteHeader(403)
				w.Write([]byte(err.Error()))
			default:
				w.WriteHeader(500)
				h.log.Error(err)
			}
		}
	})
}

// func (h *Handler) cors(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
// 		w.Header().Add("Access-Control-Allow-Headers", "*")

// 		next.ServeHTTP(w, r)
// 	})
// }

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

		userId, err := h.authHandler.services.AuthService.ParseAccessToken(headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(r.Header)

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

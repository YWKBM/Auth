package handler

import (
	"auth/handler/dto"
	"auth/services"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthorizationService
}

func newAuthHandler(authService services.AuthorizationService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.authService.CreateUser(req.Login, req.Password, req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	aToken, rToken, err := a.authService.CreateTokenPair(req.Login, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := dto.TokenPairResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (a *AuthHandler) RenewCredentials(w http.ResponseWriter, r *http.Request) {
	var req dto.RenewTokenRequest

	err := json.NewDecoder((r.Body)).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	aToken, rToken, err := a.authService.RenewToken(req.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := dto.TokenPairResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (a *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.authService.DeleteTokenPair(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (a *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	req := dto.ChangePasswordRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = a.authService.ChangePassword(userId, req.OldPassword, req.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

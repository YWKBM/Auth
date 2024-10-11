package server

import (
	"auth/server/dto"
	"auth/services"
	"encoding/json"
	"net/http"
)

type AuthServer struct {
	servs services.Services
}

func (a *AuthServer) SignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = a.servs.AuthService.CreateUser(req.Login, req.Password, req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *AuthServer) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	aToken, rToken, err := a.servs.AuthService.CreateTokenPair(req.Login, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	result := dto.TokenPairResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (a *AuthServer) RenewCredentials(w http.ResponseWriter, r *http.Request) {
	var req dto.RenewTokenRequest

	err := json.NewDecoder((r.Body)).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	aToken, rToken, err := a.servs.AuthService.RenewToken(req.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	result := dto.TokenPairResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (a *AuthServer) Logout(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthServer) ChangePassword(w http.ResponseWriter, r *http.Request) {

}

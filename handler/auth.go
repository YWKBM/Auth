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

// SignUp godoc
// @Summary      SignUp
// @Description  sign up
// @Accept       json
// @Produce      json
// @Param input body dto.SignUpRequest true "credentials"
// @Success      200
// @Failure      500
// @Router       /sign_up [post]
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

// SignIn godoc
// @Summary      SignIn
// @Description  sign in
// @Accept       json
// @Produce      json
// @Param input body dto.SignUpRequest true "credentials"
// @Success      200 	{object}	dto.TokenPairResponse
// @Failure      500
// @Router       /sign_in [post]
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

// Renew godoc
// @Summary      Renew
// @Description  renew token pair
// @Accept       json
// @Produce      json
// @Param input body dto.RenewTokenRequest true "credentials"
// @Success      200 	{object}		dto.TokenPairResponse
// @Failure      500
// @Router       /renew [post]
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

// SignOut godoc
// @Summary      SignOut
// @Description  Sign out
// @Security 	 apikey
// @Success      200
// @Failure      500
// @Router       /sign_out [post]
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

// ChangePassword godoc
// @Summary      ChangePassword
// @Description  Change password
// @Security 	 apikey
// @Param input body dto.ChangePasswordRequest true "credentials"
// @Success      200
// @Failure      500
// @Router       /change_password [post]
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

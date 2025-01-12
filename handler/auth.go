package handler

import (
	"auth/handler/dto"
	"auth/services"
	"encoding/json"
	"net/http"
	"strings"
)

type AuthHandler struct {
	services services.Services
}

func newAuthHandler(services services.Services) *AuthHandler {
	return &AuthHandler{services: services}
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
func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var req dto.SignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = req.Validate()
	if err != nil {
		return err
	}

	err = a.services.AuthService.CreateUser(req.Login, req.Password, req.Email)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthHandler) SignUpProvider(w http.ResponseWriter, r *http.Request) error {
	var req dto.ProviderSignUpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = a.services.ProviderService.RequestCreateProvider(req.FirstName, req.MiddleName, req.SecondName, req.Email, req.Phone)
	if err != nil {
		return err
	}

	return nil
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
func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var req dto.SignInRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	aToken, rToken, err := a.services.AuthService.CreateTokenPair(req.Login, req.Password)
	if err != nil {
		return err
	}

	result := dto.TokenPairResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

	return nil
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
func (a *AuthHandler) RenewCredentials(w http.ResponseWriter, r *http.Request) error {
	var req dto.RenewTokenRequest

	err := json.NewDecoder((r.Body)).Decode(&req)
	if err != nil {
		return err
	}

	aToken, rToken, err := a.services.AuthService.RenewToken(req.RefreshToken)
	if err != nil {
		return err
	}

	result := dto.TokenPairResponse{
		AccessToken:  aToken,
		RefreshToken: rToken,
	}

	resp, err := json.Marshal(result)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

	return nil
}

// SignOut godoc
// @Summary      SignOut
// @Description  Sign out
// @Security 	 apikey
// @Success      200
// @Failure      500
// @Router       /sign_out [post]
func (a *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) error {
	userId, err := getUserId(w, r)
	if err != nil {
		return err
	}

	err = a.services.AuthService.DeleteTokenPair(userId)
	if err != nil {
		return err
	}

	return nil
}

// ChangePassword godoc
// @Summary      ChangePassword
// @Description  Change password
// @Security 	 apikey
// @Param input body dto.ChangePasswordRequest true "credentials"
// @Success      200
// @Failure      500
// @Router       /change_password [post]
func (a *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	userId, err := getUserId(w, r)
	if err != nil {
		return err
	}

	req := dto.ChangePasswordRequest{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = a.services.AuthService.ChangePassword(userId, req.OldPassword, req.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthHandler) ResolveUser(w http.ResponseWriter, r *http.Request) error {
	req := dto.IdentityRequest{}
	var result dto.IdentityResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	accessToken := strings.Split(req.AuthToken, " ")
	if len(accessToken) != 2 || accessToken[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	err = a.services.AuthService.ResolveAccess(accessToken[1], req.Role)
	if err != nil {
		result = dto.IdentityResponse{
			Status: "Failed",
			Error:  err.Error(),
		}
	} else {
		result.Status = "OK"
	}

	resp, err := json.Marshal(result)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

	return nil
}

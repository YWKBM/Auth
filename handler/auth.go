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

	err = a.authService.CreateUser(req.Login, req.Password, req.Email)
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

	aToken, rToken, err := a.authService.CreateTokenPair(req.Login, req.Password)
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

	aToken, rToken, err := a.authService.RenewToken(req.RefreshToken)
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

	err = a.authService.DeleteTokenPair(userId)
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

	err = a.authService.ChangePassword(userId, req.OldPassword, req.NewPassword)
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

	err = a.authService.ResolveAccess(req.AuthToken, req.Role)
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

package services

import (
	"auth/customErrors"
	"auth/entities"
	"auth/repo"
	"auth/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	tokenTTL = 3 * time.Hour
)

type AuthService struct {
	repo       *repo.Repos
	signingKey string
}

type AccessTokenClaims struct {
	UserId   int           `json:"UserId"`
	UserRole entities.Role `json:"UserRole"`
	jwt.RegisteredClaims
}

func NewAuthService(repos *repo.Repos, signingKey string) *AuthService {
	return &AuthService{
		repo:       repos,
		signingKey: signingKey,
	}
}

func (a *AuthService) CreateUser(login, password, email string) error {
	pass := utils.GnerateHashPassword(password)
	err := a.repo.Authorization.CreateUser(login, pass, email, "USER")
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) CreateTokenPair(login, password string) (string, string, error) {
	pass := utils.GnerateHashPassword(password)
	user, err := a.repo.Authorization.GetUser(login, pass)
	if err != nil {
		return "", "", err
	}

	jti := fmt.Sprint(uuid.New())
	expiresAt := time.Now().Add(tokenTTL)

	claims := &AccessTokenClaims{
		user.Id,
		user.UserRole,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user",
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// refresh-token
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL * 2)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "user-refresh",
		ID:        jti,
	})

	st, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", "", err
	}

	sr, err := refresh.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", "", err
	}

	a.repo.Authorization.CreateToken(jti, user.Id, expiresAt)

	return st, sr, nil

}

func (a *AuthService) DeleteTokenPair(userId int) error {
	err := a.repo.Authorization.DeleteToken(userId)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ChangePassword(userId int, oldPassword, newPassword string) error {
	user, err := a.repo.Authorization.GetUserById(userId)
	if err != nil {
		return err
	}

	if user.Password != utils.GnerateHashPassword(oldPassword) {
		return errors.New("wrong login or password")
	}

	err = a.repo.Authorization.ChangePassword(user.Id, utils.GnerateHashPassword(newPassword))
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ParseAccessToken(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("unknown token claims")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return 0, err
	}

	if exp != nil && exp.Time.Unix() < time.Now().Unix() {
		return 0, errors.New("token expired")
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return 0, err
	}

	if issuer != "user" {
		return 0, errors.New("invalid token")
	}

	userId := int(claims["UserId"].(float64))

	return userId, nil
}

func (a *AuthService) RenewToken(refreshToken string) (string, string, error) {
	rToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(a.signingKey), nil
	})
	if err != nil {
		log.Fatal((err))
	}

	refreshTokenClaims, ok := rToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("unknown token claims")
	}

	issuer, err := refreshTokenClaims.GetIssuer()
	if err != nil {
		return "", "", err
	}

	if issuer != "user-refresh" {
		return "", "", errors.New("invalid token")
	}

	exp, err := refreshTokenClaims.GetExpirationTime()
	if err != nil {
		return "", "", err
	}

	if exp != nil && exp.Time.Unix() < time.Now().Unix() {
		return "", "", errors.New("token expired")
	}

	if issuer != "user-refresh" {
		return "", "", errors.New("invalid issuer")
	}

	userId, role, err := a.repo.Authorization.GetUserByTokenId(refreshTokenClaims["jti"].(string))
	if err != nil {
		return "", "", err
	}

	jti := fmt.Sprint(uuid.New())
	expiresAt := time.Now().Add(tokenTTL)

	claims := &AccessTokenClaims{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user",
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL * 2)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "user-refresh",
		ID:        jti,
	})

	st, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", "", err
	}

	sr, err := refresh.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", "", err
	}

	a.repo.Authorization.CreateToken(jti, userId, expiresAt)

	return st, sr, nil

}

func (a *AuthService) ResolveAccess(accessToken string, expectedRole string) error {
	tokenClaims, err := utils.GetTokenData(accessToken, a.signingKey)
	if err != nil {
		return err
	}

	userId := int(tokenClaims["UserId"].(float64))
	_, err = a.repo.Authorization.GetUserById(userId)
	if err != nil {
		return err
	}

	userRole, err := entities.ParseRole(tokenClaims["UserRole"].(string))
	if err != nil {
		return err
	}

	role, err := entities.ParseRole(expectedRole)

	if userRole != role {
		return &customErrors.ValidationError{Message: "forbidden for role"}
	}

	return err
}

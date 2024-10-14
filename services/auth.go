package services

import (
	"auth/repo"
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	salt       = "qweqweasddfasdfasdfqwerqwetasdg"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo *repo.Repos
}

func NewAuthService(repos *repo.Repos) *AuthService {
	return &AuthService{
		repo: repos,
	}
}

func (a *AuthService) CreateUser(login, password, email string) error {
	pass := generateHashPassword(password)
	fmt.Println(pass)
	err := a.repo.Authorization.CreateUser(login, pass, email)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) CreateProvider() (int, error) {
	return 0, nil
}

func (a *AuthService) CreateTokenPair(login, password string) (string, string, error) {
	pass := generateHashPassword(password)
	fmt.Println(pass)
	user, err := a.repo.Authorization.GetUser(login, pass)
	if err != nil {
		return "", "", err
	}

	fmt.Println(user.UserRole)
	fmt.Println(user.Email)

	jti := fmt.Sprint(uuid.New())
	expiresAt := time.Now().Add(tokenTTL)

	claims := &jwt.MapClaims{
		"UserId":    fmt.Sprintf("%v", user.Id),
		"UserRole":  user.UserRole,
		"ExpiresAt": jwt.NewNumericDate(expiresAt),
		"IssuedAt":  jwt.NewNumericDate(time.Now()),
		"Issuer":    "user",
		"jti":       jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "user-refresh",
		ID:        jti,
	})

	st, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	sr, err := refresh.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	a.repo.Authorization.CreateToken(jti, user.Id, expiresAt)
	fmt.Println(jti, user.Id, expiresAt)

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

	if user.Password != generateHashPassword(oldPassword) {
		fmt.Println("wrong password")
		return errors.New("wrong password")
	}

	err = a.repo.Authorization.ChangePassword(user.Id, generateHashPassword(newPassword))
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

		return []byte(signingKey), nil
	})
	if err != nil {
		fmt.Println(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("unknown token claims")
	}

	userIdVal := claims["UserId"]
	userId, err := strconv.Atoi(userIdVal.(string))
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (a *AuthService) RenewToken(refreshToken string) (string, string, error) {
	rToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(signingKey), nil
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
		return "", "", errors.New("invalid issuer")
	}

	userId, role, err := a.repo.Authorization.GetUserByTokenId(refreshTokenClaims["jti"].(string))
	if err != nil {
		return "", "", err
	}

	jti := fmt.Sprint(uuid.New())
	expiresAt := time.Now().Add(tokenTTL)

	claims := &jwt.MapClaims{
		"UserId":    fmt.Sprintf("%v", userId),
		"UserRole":  role,
		"ExpiresAt": jwt.NewNumericDate(expiresAt),
		"IssuedAt":  jwt.NewNumericDate(time.Now()),
		"Issuer":    "user",
		"jti":       jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "user-refresh",
		ID:        jti,
	})

	st, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	sr, err := refresh.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	a.repo.Authorization.CreateToken(jti, userId, expiresAt)

	return st, sr, nil

}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

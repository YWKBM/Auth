package services

import (
	"auth/entities"
	"auth/repo"
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	salt       = "qweqweasddfasdfasdfqwerqwetasdg"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type accessTokenClaims struct {
	UserId   int
	UserRole entities.Role
	jwt.RegisteredClaims
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
}

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
	user, err := a.repo.Authorization.GetUser(login, generateHashPassword(password))
	if err != nil {
		return "", "", err
	}

	jti := fmt.Sprint(uuid.New())
	expiresAt := time.Now().Add(tokenTTL)

	a.repo.Authorization.CreateToken(jti, user.Id, expiresAt)

	claims := accessTokenClaims{
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

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "user-refresh",
		ID:        jti,
	})

	st, err := token.SignedString(signingKey)
	if err != nil {
		return "", "", err
	}

	sr, err := refresh.SignedString(signingKey)
	if err != nil {
		return "", "", err
	}

	return st, sr, nil

}

func (a *AuthService) ParseAccessToken(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	claims, ok := token.Claims.(*accessTokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *accessTokenClaims")
	}

	return claims.UserId, nil
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

	refreshTokenClaims, ok := rToken.Claims.(*refreshTokenClaims)
	if !ok {
		return "", "", errors.New("token claims are not of type *refreshTokenClaims")
	}

	userId, role, err := a.repo.Authorization.GetUserByTokenId(refreshTokenClaims.ID)
	if err != nil {
		return "", "", err
	}

	jti := fmt.Sprint(uuid.New())
	expiresAt := time.Now().Add(tokenTTL)

	a.repo.Authorization.CreateToken(jti, userId, expiresAt)

	claims := accessTokenClaims{
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
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "user-refresh",
		ID:        jti,
	})

	st, err := token.SignedString(signingKey)
	if err != nil {
		return "", "", err
	}

	sr, err := refresh.SignedString(signingKey)
	if err != nil {
		return "", "", err
	}

	return st, sr, nil

}

func generateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

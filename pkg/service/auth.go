package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"time"
	"todo"
	"todo/pkg/repository"

	"github.com/golang-jwt/jwt"
)

var salt = os.Getenv("HASH_SALT");
const tokenTTL = 12 * time.Hour

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user) 
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	//get user from DB
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt: time.Now().Unix(),},
		user.ID,
	})

	return token.SignedString([]byte(salt))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token ,err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(accessToken *jwt.Token) (interface{}, error) {
		if _, ok := accessToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(salt), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok:=token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token is expired")
	}

	return claims.UserId, nil
}

func  generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return  fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
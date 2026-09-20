package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService() *AuthService {
	return &AuthService{jwtSecret: []byte(os.Getenv("JWT_SECRET"))}
}

func (a *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (a *AuthService) CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (a *AuthService) GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func (a *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token tidak valid")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("token tidak valid")
	}

	return uint(userIDFloat), nil
}
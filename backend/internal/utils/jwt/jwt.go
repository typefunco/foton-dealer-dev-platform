package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims представляет claims для JWT токена
type JWTClaims struct {
	Login   string `json:"login"`
	IsAdmin bool   `json:"is_admin"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

var (
	secretKey = []byte("secret-key")
	ttl       = time.Hour * 24 * 7 // Увеличиваем время жизни токена до 7 дней
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateJWT(login string, isAdmin bool, role string) (string, error) {
	claims := JWTClaims{
		Login:   login,
		IsAdmin: isAdmin,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateJWTLegacy поддерживает старый метод для обратной совместимости
func (s *Service) ValidateJWTLegacy(tokenString string) error {
	_, err := s.ValidateJWT(tokenString)
	return err
}

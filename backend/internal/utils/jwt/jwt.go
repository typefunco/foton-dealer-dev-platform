package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

//type JWTRepository interface {
//	ValidateJWT(ctx context.Context, jwt string) error
//	GenerateJWT(ctx context.Context, user model.User) (string, error)
//}

var (
	secretKey = []byte("secret-key")
	ttl       = time.Hour * 24 * 60
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateJWT(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"login": login,
			"exp":   time.Now().Add(ttl).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

package auth

import (
	"fmt"
	"time"

	"consistent_1/Domain"

	"github.com/dgrijalva/jwt-go"
)


type JWTService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDFromToken(token string) (string, error)
}

type jwtCustomClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}


func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
		issuer:    "consistent_1",
	}
}


func (service *jwtService) GenerateToken(userID string) (string, error) {
	claims := &jwtCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), 
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(service.secretKey))
}


func (service *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}
func (service *jwtService) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := service.ValidateToken(tokenString)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return "", domain.ErrInvalidToken
	}

	return claims.UserID, nil
}
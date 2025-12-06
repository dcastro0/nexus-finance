package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey string
	issuer    string
}

func NewTokenService(secretKey, issuer string) *TokenService {
	return &TokenService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (s *TokenService) GenerateToken(accountID string) (string, error) {
	claims := &Claims{
		AccountID: accountID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

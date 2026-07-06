package security

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint
	Role   string
	Iat    int64
	Exp    int64

	jwt.RegisteredClaims
}

type tokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string, secret string, expirationHours int) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(expirationHours) * time.Hour)

	claims := tokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	parsedClaims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if parsedClaims.Subject == "" {
		return nil, fmt.Errorf("invalid token subject")
	}

	parsedUserID, err := strconv.ParseUint(parsedClaims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid token subject")
	}

	claims := &Claims{
		UserID: uint(parsedUserID),
		Role:   parsedClaims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   parsedClaims.Subject,
			IssuedAt:  parsedClaims.IssuedAt,
			ExpiresAt: parsedClaims.ExpiresAt,
		},
	}

	if parsedClaims.IssuedAt != nil {
		claims.Iat = parsedClaims.IssuedAt.Unix()
	}
	if parsedClaims.ExpiresAt != nil {
		claims.Exp = parsedClaims.ExpiresAt.Unix()
	}

	return claims, nil
}

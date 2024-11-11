package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	Claims
	jwt.StandardClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var jwtSecret = []byte("super-secret-key") // In production, use environment variable

func GenerateToken(userID int64) (string, error) {
	claims := &JWTClaim{
		Claims: Claims{UserID: userID},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 24 hours for access token
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Changed from ES256 to HS256 to match the validation method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userID int64) (string, error) {
	claims := &JWTClaim{
		Claims: Claims{UserID: userID},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days for refresh token
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			// Check if the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return &claims.Claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateTokenPair(userID int64) (*TokenPair, error) {
	// Generate access token
	accessToken, err := GenerateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	// Generate refresh token
	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateRefreshToken(tokenString string) (int64, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return 0, fmt.Errorf("invalid refresh token: %v", err)
	}

	return claims.UserID, nil
}
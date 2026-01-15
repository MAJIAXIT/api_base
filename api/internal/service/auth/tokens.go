package auth

import (
	"errors"
	"time"

	"github.com/MAJIAXIT/api_base/api/pkg/logger"
	"github.com/MAJIAXIT/api_base/api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	Type   string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

func (s *service) GenerateTokens(
	tx *gorm.DB,
	userID uint,
	login string,
	userAgent string,
	ip string) (
	string, string, error) {

	accessToken, err := s.GenerateToken(userID, login, AccessToken)
	if err != nil {
		return "", "", utils.WrapError(err)
	}

	refreshToken, err := s.GenerateToken(userID, login, RefreshToken)
	if err != nil {
		return "", "", utils.WrapError(err)
	}

	_, err = s.CreateSessionWithToken(tx, refreshToken, userID, userAgent, ip)
	if err != nil {
		return "", "", utils.WrapError(err)
	}

	return accessToken, refreshToken, nil
}

func (s *service) GenerateToken(
	userID uint,
	login string,
	tokenType TokenType) (
	string, error) {

	// Set expiration time based on token type
	var expirationTime time.Duration
	if tokenType == AccessToken {
		expirationTime = s.JWTConfig.AccessTokenExpireDuration
	} else {
		expirationTime = s.JWTConfig.RefreshTokenExpireDuration
	}

	// Create claims
	claims := TokenClaims{
		UserID: userID,
		Login:  login,
		Type:   string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(s.JWTConfig.Secret))
	if err != nil {
		return "", logger.WrapError(err)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (s *service) ValidateToken(
	tokenString string, expectedType TokenType) (
	*TokenClaims, error) {

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, utils.NewUnauthorized("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.JWTConfig.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, utils.NewUnauthorized("token is expired")
		}
		return nil, logger.WrapError(err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		// Validate token type
		if TokenType(claims.Type) != expectedType {
			return nil, utils.NewUnauthorized("invalid token type")
		}
		return claims, nil
	}

	return nil, utils.NewUnauthorized("invalid token")
}

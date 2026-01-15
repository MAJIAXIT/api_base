package auth

import (
	"github.com/MAJIAXIT/api_base/api/config"
	"github.com/MAJIAXIT/api_base/api/internal/dto/auth"
	"github.com/MAJIAXIT/api_base/api/internal/models/session"
	"github.com/MAJIAXIT/api_base/api/internal/models/user"
	"github.com/MAJIAXIT/api_base/api/internal/service/users"
	"gorm.io/gorm"
)

type Service interface {
	Authenticate(tx *gorm.DB, req *auth.LoginRequest) (*user.User, error)
	CreateSessionWithToken(tx *gorm.DB, token string, userID uint, userAgent string, ip string) (*session.Session, error)
	ValidateSessionByToken(tx *gorm.DB, token string) (*session.Session, error)
	DeleteSessionByToken(tx *gorm.DB, token string) error
	DeleteSessionsByUserID(tx *gorm.DB, userID uint) error
	CleanupExpiredSessions(tx *gorm.DB) error
	GenerateTokens(tx *gorm.DB, userID uint, login string, userAgent string, ip string) (string, string, error)
	GenerateToken(userID uint, login string, tokenType TokenType) (string, error)
	ValidateToken(tokenString string, expectedType TokenType) (*TokenClaims, error)
}

type service struct {
	usersService users.Service
	JWTConfig    config.JWTConfig
}

func New(jwtConfig *config.JWTConfig, usersService users.Service) Service {
	return &service{
		JWTConfig:    *jwtConfig,
		usersService: usersService,
	}
}

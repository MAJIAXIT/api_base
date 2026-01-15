package auth

import (
	"errors"
	"time"

	"github.com/MAJIAXIT/api_base/api/internal/models/session"
	"github.com/MAJIAXIT/api_base/api/pkg/logger"
	"github.com/MAJIAXIT/api_base/api/pkg/utils"
	"gorm.io/gorm"
)

func (s *service) CreateSessionWithToken(
	tx *gorm.DB,
	token string,
	userID uint,
	userAgent string,
	ip string) (
	*session.Session, error) {

	session := &session.Session{
		UserID:    userID,
		Token:     token,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := tx.Create(session).Error; err != nil {
		return nil, logger.WrapError(err)
	}

	return session, nil
}

func (s *service) ValidateSessionByToken(
	tx *gorm.DB, token string) (
	*session.Session, error) {

	var session session.Session
	if err := tx.
		Where("token = ? AND expires_at > ?",
			token, time.Now()).
		First(&session).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewUnauthorized("session invalid or expired")
		}
		return nil, logger.WrapError(err)
	}
	return &session, nil
}

func (s *service) DeleteSessionByToken(
	tx *gorm.DB, token string) error {

	if err := tx.
		Where("token = ?", token).
		Delete(&session.Session{}).
		Error; err != nil {
		return logger.WrapError(err)
	}
	return nil
}

func (s *service) DeleteSessionsByUserID(
	tx *gorm.DB, userID uint) error {

	if err := tx.
		Where("user_id = ?", userID).
		Delete(&session.Session{}).
		Error; err != nil {
		return logger.WrapError(err)
	}
	return nil
}

func (s *service) CleanupExpiredSessions(
	tx *gorm.DB) error {

	if err := tx.
		Where("expires_at < ?", time.Now()).
		Delete(&session.Session{}).
		Error; err != nil {
		return logger.WrapError(err)
	}
	return nil
}

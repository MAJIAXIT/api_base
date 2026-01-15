package auth

import (
	"strings"

	"github.com/MAJIAXIT/api_base/api/internal/dto/auth"
	"github.com/MAJIAXIT/api_base/api/internal/models/user"
	"github.com/MAJIAXIT/api_base/api/pkg/logger"
	"github.com/MAJIAXIT/api_base/api/pkg/utils"
	"gorm.io/gorm"
)

func (s *service) Authenticate(
	tx *gorm.DB, req *auth.LoginRequest) (
	*user.User, error) {

	user, err := s.usersService.GetUserByLoginOrEmail(tx, strings.TrimSpace(req.Login))
	if err != nil {
		return nil, utils.WrapError(err)
	}

	match, err := user.ComparePasswords(req.Password)
	if err != nil {
		return nil, logger.WrapError(err)
	}
	if !match {
		return nil, utils.NewUnauthorized("Invalid credentials")
	}

	return user, nil
}

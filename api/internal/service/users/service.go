package users

import (
	auth_dto "github.com/MAJIAXIT/api_base/api/internal/dto/auth"
	users_dto "github.com/MAJIAXIT/api_base/api/internal/dto/users"
	"github.com/MAJIAXIT/api_base/api/internal/models/user"
	"gorm.io/gorm"
)

type Service interface {
	GetUserByLoginOrEmail(tx *gorm.DB, login string) (*user.User, error)
	GetUserByID(tx *gorm.DB, id uint) (*user.User, error)
	UserBeforeCreateExistsCheck(tx *gorm.DB, login string) (bool, error)
	CreateUser(tx *gorm.DB, req *auth_dto.SignupRequest) (*user.User, error)
	UpdateUser(tx *gorm.DB, userID uint, req *users_dto.UpdateUserRequest) (*user.User, error)
}

type service struct {
}

func New() Service {
	return &service{}
}

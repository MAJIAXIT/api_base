package users

import (
	"errors"

	auth_dto "github.com/MAJIAXIT/api_base/api/internal/dto/auth"
	users_dto "github.com/MAJIAXIT/api_base/api/internal/dto/users"
	"github.com/MAJIAXIT/api_base/api/internal/models/user"
	"github.com/MAJIAXIT/api_base/api/pkg/logger"
	"github.com/MAJIAXIT/api_base/api/pkg/utils"
	"gorm.io/gorm"
)

func (s *service) GetUserByLoginOrEmail(
	tx *gorm.DB, login string) (
	*user.User, error) {

	var user user.User
	if err := tx.
		// Where("login = ? OR email = ?", login, email).
		Where("login = ?", login).
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewUnauthorized("Invalid credentials")
		}
		return nil, logger.WrapError(err)
	}

	return &user, nil
}

func (s *service) GetUserByID(
	tx *gorm.DB, id uint) (
	*user.User, error) {

	var user user.User
	if err := tx.
		Where("id = ?", id).
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("User with id: %d not found", id)
		}
		return nil, logger.WrapError(err)
	}

	return &user, nil
}

func (s *service) UserBeforeCreateExistsCheck(
	tx *gorm.DB, login string) (
	bool, error) {

	var count int64
	if err := tx.
		Model(&user.User{}).
		Where("login = ?", login).
		Count(&count).
		Error; err != nil {
		return false, logger.WrapError(err)
	}
	if count > 0 {
		return true, utils.NewConflict("Пользователь с таким login уже существует")
	}

	return count > 0, nil
}

func (s *service) CreateUser(
	tx *gorm.DB, req *auth_dto.SignupRequest) (
	*user.User, error) {

	usr := user.User{
		EncrPassword: req.Password,
		Login:        req.Login,
	}
	if err := usr.EncryptPassword(usr.EncrPassword); err != nil {
		return nil, logger.WrapError(err)
	}

	if err := tx.Create(&usr).Error; err != nil {
		return nil, logger.WrapError(err)
	}

	return &usr, nil
}

func (s *service) UpdateUser(
	tx *gorm.DB,
	userID uint,
	req *users_dto.UpdateUserRequest) (
	*user.User, error) {

	usr, err := s.GetUserByID(tx, userID)
	if err != nil {
		return nil, utils.WrapError(err)
	}

	var count int64
	if req.Login != "" && req.Login != usr.Login {
		if err := tx.
			Model(&user.User{}).
			Where("login = ?", req.Login).
			Count(&count).
			Error; err != nil {
			return nil, logger.WrapError(err)
		}
		if count > 0 {
			return nil, utils.NewConflict("Пользователь с таким login уже существует")
		}
		usr.Login = req.Login
	}

	if err := tx.Save(usr).Error; err != nil {
		return nil, logger.WrapErrMsg(err, "Failed to save user")
	}

	usr.EncrPassword = ""

	return usr, nil
}

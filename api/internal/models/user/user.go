package user

import (
	"github.com/MAJIAXIT/projname/api/internal/models/base"
	"github.com/MAJIAXIT/projname/api/pkg/logger"
	"github.com/MAJIAXIT/projname/api/pkg/password"
)

type User struct {
	base.BaseModel
	Login        string `json:"login" gorm:"uniqueIndex:idx_user_email"`
	EncrPassword string `json:"encr_password,omitempty"`
}

// PasswordHolder interface implementation
func (u *User) GetEncrPassword() string {
	return u.EncrPassword
}

func (u *User) SetEncrPassword(encrPassword string) {
	u.EncrPassword = encrPassword
}

func (u *User) EncryptPassword(plainPassword string) error {
	return logger.WrapError(password.Encrypt(u, plainPassword))
}

func (u *User) ComparePasswords(inputPassword string) (match bool, err error) {
	return password.Compare(u, inputPassword)
}

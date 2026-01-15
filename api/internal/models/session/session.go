package session

import (
	"time"

	"github.com/MAJIAXIT/api_base/api/internal/models/base"
	"github.com/MAJIAXIT/api_base/api/internal/models/user"
)

type Session struct {
	base.BaseModel

	UserID    uint       `json:"user_id" gorm:"not null;index:idx_session_user;constraint:OnDelete:CASCADE"`
	User      *user.User `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Token     string     `json:"token" gorm:"uniqueIndex"`
	UserAgent string     `json:"user_agent"`
	IP        string     `json:"ip"`
	ExpiresAt time.Time  `json:"expires_at"`

	_ struct{} `gorm:"index:idx_session_user_expires,composite:user_id,expires_at"`
	_ struct{} `gorm:"index:idx_session_token_expires,composite:token,expires_at"`
}

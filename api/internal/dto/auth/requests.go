package auth

type LoginRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32,ascii"`
	Password string `json:"password" binding:"required,ascii"`
}

type SignupRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32,ascii"`
	Password string `json:"password" binding:"required,min=8,max=128,ascii"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,jwt"`
}

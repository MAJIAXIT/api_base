package users

type UpdateUserRequest struct {
	Login string `json:"login,omitempty" binding:"required,min=3,max=32,ascii"`
}

package dto

type UserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	FullName      string `json:"full_name"`
	AvatarUrl     string `json:"avatar_url"`
	Bio           string `json:"bio"`
	IsActive      bool   `json:"is_active"`
	EmailVerified bool   `json:"email_verified"`
}

type UpdateProfileRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	FullName  string `json:"full_name" validate:"required"`
	AvatarUrl string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

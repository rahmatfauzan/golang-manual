package dto

type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FullName  string `json:"full_name" validate:"required"`
	AvatarUrl string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
}

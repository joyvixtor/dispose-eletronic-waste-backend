package userauth

type RegisterRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	UFPEEmail string `json:"ufpe_email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Workplace string `json:"workplace" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=admin bolsista servidor"`
}

type LoginRequest struct {
	UFPEEmail string `json:"ufpe_email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	FullName  string `json:"full_name"`
	UFPEEmail string `json:"ufpe_email"`
	Workplace string `json:"workplace"`
	Role      string `json:"role"`
}

package dto

type DummyLoginRequest struct {
	Role string `binding:"required,oneof=admin user" json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Email    string `binding:"required,email"            json:"email"`
	Password string `binding:"required,min=6"            json:"password"`
}

type LoginRequest struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required"       json:"password"`
}

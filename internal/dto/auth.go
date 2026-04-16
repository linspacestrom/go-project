package dto

import "time"

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Email      string     `binding:"required,email" json:"email"`
	Password   string     `binding:"required,min=8" json:"password"`
	FullName   string     `binding:"required" json:"full_name"`
	BirthDate  *time.Time `json:"birth_date"`
	University string     `json:"university"`
	Course     *int       `json:"course"`
	DegreeType string     `json:"degree_type"`
	Role       string     `json:"role"`
}

type LoginRequest struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required"       json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `binding:"required" json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `binding:"required" json:"refresh_token"`
}

type RegisterMentorRequest struct {
	Email       string  `binding:"required,email" json:"email"`
	Password    string  `binding:"required,min=8" json:"password"`
	FullName    string  `binding:"required" json:"full_name"`
	Description *string `json:"description"`
	Title       *string `json:"title"`
	CityID      string  `binding:"required,uuid" json:"city_id"`
}

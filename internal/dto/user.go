package dto

import "time"

type UserResponse struct {
	ID        string     `json:"id"`
	FullName  string     `json:"full_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CityID    *string    `json:"city_id,omitempty"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UpdateMeRequest struct {
	FullName  *string    `json:"full_name"`
	BirthDate *time.Time `json:"birth_date"`
}

type UpdateUserCityRequest struct {
	CityID string `binding:"required,uuid" json:"city_id"`
}

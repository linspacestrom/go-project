package dto

import "time"

type SlotFilter struct {
	Pagination
	CityID    string     `form:"city_id"`
	HubID     string     `form:"hub_id"`
	RoomID    string     `form:"room_id"`
	MentorID  string     `form:"mentor_id"`
	Status    string     `form:"status"`
	Type      string     `form:"type"`
	StartFrom *time.Time `form:"start_from" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTo     *time.Time `form:"end_to" time_format:"2006-01-02T15:04:05Z07:00"`
	SortBy    string     `form:"sort_by"`
	Order     string     `form:"order"`
}

type CreateSlotRequest struct {
	RoomID     *string   `json:"room_id"`
	CityID     string    `binding:"required,uuid" json:"city_id"`
	Type       string    `binding:"required,oneof=online offline" json:"type"`
	StartAt    time.Time `binding:"required" json:"start_at"`
	EndAt      time.Time `binding:"required" json:"end_at"`
	MeetingURL *string   `json:"meeting_url"`
	Status     string    `json:"status"`
	Capacity   *int      `json:"capacity"`
}

type UpdateSlotRequest struct {
	RoomID     *string    `json:"room_id"`
	StartAt    *time.Time `json:"start_at"`
	EndAt      *time.Time `json:"end_at"`
	MeetingURL *string    `json:"meeting_url"`
	Status     *string    `json:"status"`
	Capacity   *int       `json:"capacity"`
}

type SlotResponse struct {
	ID         string    `json:"id"`
	MentorID   string    `json:"mentor_id"`
	RoomID     *string   `json:"room_id,omitempty"`
	CityID     string    `json:"city_id"`
	Type       string    `json:"type"`
	StartAt    time.Time `json:"start_at"`
	EndAt      time.Time `json:"end_at"`
	MeetingURL *string   `json:"meeting_url,omitempty"`
	Status     string    `json:"status"`
	Capacity   *int      `json:"capacity,omitempty"`
}

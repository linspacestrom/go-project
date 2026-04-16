package dto

import "time"

type BookingFilter struct {
	Pagination
	Status   string `form:"status"`
	CityID   string `form:"city_id"`
	RoomID   string `form:"room_id"`
	SlotID   string `form:"slot_id"`
	MentorID string `form:"mentor_id"`
	SortBy   string `form:"sort_by"`
	Order    string `form:"order"`
}

type CreateBookingRequest struct {
	SlotID         *string   `json:"slot_id"`
	RoomID         *string   `json:"room_id"`
	BookingType    string    `binding:"required,oneof=room_only room_with_mentor mentor_call event_seat" json:"booking_type"`
	StartAt        time.Time `binding:"required" json:"start_at"`
	EndAt          time.Time `binding:"required" json:"end_at"`
	MeetingURL     *string   `json:"meeting_url"`
	SeatNumber     *int      `json:"seat_number"`
	IdempotencyKey *string   `json:"idempotency_key"`
}

type BookingResponse struct {
	ID          string    `json:"id"`
	SlotID      *string   `json:"slot_id,omitempty"`
	RoomID      *string   `json:"room_id,omitempty"`
	UserID      string    `json:"user_id"`
	BookingType string    `json:"booking_type"`
	Status      string    `json:"status"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	MeetingURL  *string   `json:"meeting_url,omitempty"`
	SeatNumber  *int      `json:"seat_number,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

package dto

type CreateCityRequest struct {
	Name     string `binding:"required" json:"name"`
	IsActive *bool  `json:"is_active"`
}

type UpdateCityRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

type CityResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type CreateHubRequest struct {
	CityID   string `binding:"required,uuid" json:"city_id"`
	Name     string `binding:"required" json:"name"`
	Address  string `binding:"required" json:"address"`
	IsActive *bool  `json:"is_active"`
}

type UpdateHubRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	IsActive *bool   `json:"is_active"`
}

type HubResponse struct {
	ID       string `json:"id"`
	CityID   string `json:"city_id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`
}

type CreateRoomRequest struct {
	HubID       string  `binding:"required,uuid" json:"hub_id"`
	Name        string  `binding:"required" json:"name"`
	Description *string `json:"description"`
	RoomType    *string `json:"room_type"`
	Capacity    int     `binding:"required,min=1" json:"capacity"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateRoomRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	RoomType    *string `json:"room_type"`
	Capacity    *int    `json:"capacity"`
	IsActive    *bool   `json:"is_active"`
}

type RoomResponse struct {
	ID          string  `json:"id"`
	HubID       string  `json:"hub_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	RoomType    *string `json:"room_type,omitempty"`
	Capacity    int     `json:"capacity"`
	IsActive    bool    `json:"is_active"`
}

type BusinessAnalyticsResponse struct {
	MentorApprovals                int64 `json:"mentor_approvals"`
	MentorRejections               int64 `json:"mentor_rejections"`
	ActiveBookings                 int64 `json:"active_bookings"`
	MentorRequestsBySkill          int64 `json:"mentor_requests_by_skill"`
	MentorRequestsByCategory       int64 `json:"mentor_requests_by_category"`
	UserCancelsAfterMentorApproval int64 `json:"user_cancels_after_mentor_approval"`
}

type TechnicalAnalyticsResponse struct {
	BookingConflictCount int64 `json:"booking_conflict_count"`
	FailedOutboxEvents   int64 `json:"failed_outbox_events"`
	OutboxLagCount       int64 `json:"outbox_lag_count"`
}

package dto

type MentorRequestFilter struct {
	Pagination
	Status   string `form:"status"`
	SkillID  string `form:"skill_id"`
	CityID   string `form:"city_id"`
	MentorID string `form:"mentor_id"`
}

type CreateMentorRequestRequest struct {
	SlotID      *string `json:"slot_id"`
	MentorID    *string `json:"mentor_id"`
	SkillID     *string `json:"skill_id"`
	RequestType string  `binding:"required,oneof=category skill direct_mentor other" json:"request_type"`
	Comment     *string `json:"comment"`
}

type MentorRequestResponse struct {
	ID          string  `json:"id"`
	SlotID      *string `json:"slot_id,omitempty"`
	MentorID    *string `json:"mentor_id,omitempty"`
	MenteeID    string  `json:"mentee_id"`
	CityID      string  `json:"city_id"`
	SkillID     *string `json:"skill_id,omitempty"`
	RequestType string  `json:"request_type"`
	Status      string  `json:"status"`
	Comment     *string `json:"comment,omitempty"`
}

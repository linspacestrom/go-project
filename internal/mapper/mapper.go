package mapper

import (
	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/dto"
)

func ToUserResponse(u *domain.User) dto.UserResponse {
	res := dto.UserResponse{
		ID:        u.ID.String(),
		FullName:  u.FullName,
		BirthDate: u.BirthDate,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.CityID != nil {
		cityID := u.CityID.String()
		res.CityID = &cityID
	}

	return res
}

func ToCityResponse(c *domain.City) dto.CityResponse {
	return dto.CityResponse{ID: c.ID.String(), Name: c.Name, IsActive: c.IsActive}
}

func ToHubResponse(h *domain.Hub) dto.HubResponse {
	return dto.HubResponse{ID: h.ID.String(), CityID: h.CityID.String(), Name: h.Name, Address: h.Address, IsActive: h.IsActive}
}

func ToRoomResponse(r *domain.Room) dto.RoomResponse {
	return dto.RoomResponse{
		ID:          r.ID.String(),
		HubID:       r.HubID.String(),
		Name:        r.Name,
		Description: r.Description,
		RoomType:    r.RoomType,
		Capacity:    r.Capacity,
		IsActive:    r.IsActive,
	}
}

func ToSlotResponse(s *domain.Slot) dto.SlotResponse {
	res := dto.SlotResponse{
		ID:         s.ID.String(),
		MentorID:   s.MentorID.String(),
		CityID:     s.CityID.String(),
		Type:       s.Type,
		StartAt:    s.StartAt,
		EndAt:      s.EndAt,
		MeetingURL: s.MeetingURL,
		Status:     s.Status,
		Capacity:   s.Capacity,
	}
	if s.RoomID != nil {
		roomID := s.RoomID.String()
		res.RoomID = &roomID
	}

	return res
}

func ToBookingResponse(b *domain.Booking) dto.BookingResponse {
	res := dto.BookingResponse{
		ID:          b.ID.String(),
		UserID:      b.UserID.String(),
		BookingType: b.BookingType,
		Status:      b.Status,
		StartAt:     b.StartAt,
		EndAt:       b.EndAt,
		MeetingURL:  b.MeetingURL,
		SeatNumber:  b.SeatNumber,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
	if b.SlotID != nil {
		slotID := b.SlotID.String()
		res.SlotID = &slotID
	}
	if b.RoomID != nil {
		roomID := b.RoomID.String()
		res.RoomID = &roomID
	}

	return res
}

func ToMentorRequestResponse(r *domain.MentorRequest) dto.MentorRequestResponse {
	res := dto.MentorRequestResponse{
		ID:          r.ID.String(),
		MenteeID:    r.MenteeID.String(),
		CityID:      r.CityID.String(),
		RequestType: r.RequestType,
		Status:      r.Status,
		Comment:     r.Comment,
	}
	if r.SlotID != nil {
		slotID := r.SlotID.String()
		res.SlotID = &slotID
	}
	if r.MentorID != nil {
		mentorID := r.MentorID.String()
		res.MentorID = &mentorID
	}
	if r.SkillID != nil {
		skillID := r.SkillID.String()
		res.SkillID = &skillID
	}

	return res
}

func ParseUUIDPtr(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	parsed, err := uuid.Parse(*value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

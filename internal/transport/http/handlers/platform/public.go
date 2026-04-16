package platform

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/dto"
	"github.com/linspacestrom/go-project/internal/mapper"
	"github.com/linspacestrom/go-project/internal/repository/postgres"
	"github.com/linspacestrom/go-project/internal/server/middleware"
)

func (h *Handler) GetMe(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}

	user, err := h.s.GetMe(c.Request.Context(), userID)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToUserResponse(user))
}

func (h *Handler) UpdateMe(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}

	var req dto.UpdateMeRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	user, err := h.s.UpdateMe(c.Request.Context(), userID, req.FullName, req.BirthDate)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToUserResponse(user))
}

func (h *Handler) ListCities(c *gin.Context) {
	var p dto.Pagination
	_ = c.ShouldBindQuery(&p)
	p = p.Normalize(20, 100)

	cities, total, err := h.s.ListCities(c.Request.Context(), p.Limit, p.Offset)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, toCityListResponse(cities, total, p.Limit, p.Offset))
}

func (h *Handler) ListHubsByCity(c *gin.Context) {
	cityID, err := uuid.Parse(c.Param("cityId"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var p dto.Pagination
	_ = c.ShouldBindQuery(&p)
	p = p.Normalize(20, 100)

	hubs, total, err := h.s.ListHubsByCity(c.Request.Context(), cityID, p.Limit, p.Offset)
	if err != nil {
		h.mapError(c, err)
		return
	}

	items := make([]dto.HubResponse, 0, len(hubs))
	for _, hub := range hubs {
		hubCopy := hub
		items = append(items, mapper.ToHubResponse(&hubCopy))
	}
	c.JSON(http.StatusOK, dto.ListResponse[dto.HubResponse]{Items: items, TotalCount: total, Limit: p.Limit, Offset: p.Offset})
}

func (h *Handler) ListRoomsByCity(c *gin.Context) {
	cityID, err := uuid.Parse(c.Param("cityId"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var p dto.Pagination
	_ = c.ShouldBindQuery(&p)
	p = p.Normalize(20, 100)

	rooms, total, err := h.s.ListRoomsByCity(c.Request.Context(), cityID, p.Limit, p.Offset)
	if err != nil {
		h.mapError(c, err)
		return
	}

	items := make([]dto.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomCopy := room
		items = append(items, mapper.ToRoomResponse(&roomCopy))
	}
	c.JSON(http.StatusOK, dto.ListResponse[dto.RoomResponse]{Items: items, TotalCount: total, Limit: p.Limit, Offset: p.Offset})
}

func (h *Handler) ListSlots(c *gin.Context) {
	var req dto.SlotFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	req.Pagination = req.Pagination.Normalize(20, 100)

	cityID, err := ptrOrNilUUID(req.CityID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	hubID, err := ptrOrNilUUID(req.HubID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	roomID, err := ptrOrNilUUID(req.RoomID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	mentorID, err := ptrOrNilUUID(req.MentorID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}
	var slotType *string
	if req.Type != "" {
		slotType = &req.Type
	}

	slots, total, err := h.s.ListSlots(c.Request.Context(), postgres.SlotFilter{
		CityID:    cityID,
		HubID:     hubID,
		RoomID:    roomID,
		MentorID:  mentorID,
		Status:    status,
		Type:      slotType,
		StartFrom: req.StartFrom,
		EndTo:     req.EndTo,
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
	if err != nil {
		h.mapError(c, err)
		return
	}

	items := make([]dto.SlotResponse, 0, len(slots))
	for _, slot := range slots {
		slotCopy := slot
		items = append(items, mapper.ToSlotResponse(&slotCopy))
	}

	c.JSON(http.StatusOK, dto.ListResponse[dto.SlotResponse]{Items: items, TotalCount: total, Limit: req.Limit, Offset: req.Offset})
}

func (h *Handler) GetSlotByID(c *gin.Context) {
	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	slot, err := h.s.GetSlotByID(c.Request.Context(), slotID)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToSlotResponse(slot))
}

func (h *Handler) CreateBooking(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	var req dto.CreateBookingRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	slotID, err := mapper.ParseUUIDPtr(req.SlotID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	roomID, err := mapper.ParseUUIDPtr(req.RoomID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	booking, err := h.s.CreateBooking(c.Request.Context(), userID, &domain.Booking{
		SlotID:         slotID,
		RoomID:         roomID,
		BookingType:    req.BookingType,
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		MeetingURL:     req.MeetingURL,
		SeatNumber:     req.SeatNumber,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.ToBookingResponse(booking))
}

func (h *Handler) CancelBooking(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	canceled, err := h.s.CancelBooking(c.Request.Context(), bookingID, userID, false)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToBookingResponse(canceled))
}

func (h *Handler) ListBookings(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	var req dto.BookingFilter
	if err = c.ShouldBindQuery(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	req.Pagination = req.Pagination.Normalize(20, 100)

	roomID, err := ptrOrNilUUID(req.RoomID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	slotID, err := ptrOrNilUUID(req.SlotID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	cityID, err := ptrOrNilUUID(req.CityID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	var status *string
	if req.Status != "" {
		status = &req.Status
	}

	bookings, total, err := h.s.ListBookings(c.Request.Context(), postgres.BookingFilter{
		UserID: &userID,
		Status: status,
		CityID: cityID,
		RoomID: roomID,
		SlotID: slotID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		h.mapError(c, err)
		return
	}

	items := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		copy := b
		items = append(items, mapper.ToBookingResponse(&copy))
	}

	c.JSON(http.StatusOK, dto.ListResponse[dto.BookingResponse]{Items: items, TotalCount: total, Limit: req.Limit, Offset: req.Offset})
}

func (h *Handler) CreateMentorRequest(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	user, err := h.s.GetMe(c.Request.Context(), userID)
	if err != nil {
		h.mapError(c, err)
		return
	}
	if user.CityID == nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	var req dto.CreateMentorRequestRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	slotID, err := mapper.ParseUUIDPtr(req.SlotID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	mentorID, err := mapper.ParseUUIDPtr(req.MentorID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	skillID, err := mapper.ParseUUIDPtr(req.SkillID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	created, err := h.s.CreateMentorRequest(c.Request.Context(), &domain.MentorRequest{
		SlotID:      slotID,
		MentorID:    mentorID,
		MenteeID:    userID,
		CityID:      *user.CityID,
		SkillID:     skillID,
		RequestType: req.RequestType,
		Comment:     req.Comment,
	})
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.ToMentorRequestResponse(created))
}

func (h *Handler) ListMentorRequests(c *gin.Context) {
	userID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	var req dto.MentorRequestFilter
	if err = c.ShouldBindQuery(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	req.Pagination = req.Pagination.Normalize(20, 100)

	skillID, err := ptrOrNilUUID(req.SkillID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	cityID, err := ptrOrNilUUID(req.CityID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	mentorID, err := ptrOrNilUUID(req.MentorID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	role, _ := middleware.GetRole(c)
	filter := postgres.MentorRequestFilter{SkillID: skillID, CityID: cityID, MentorID: mentorID, Limit: req.Limit, Offset: req.Offset}
	if role == domain.RoleStudent {
		filter.MenteeID = &userID
	}
	if role == domain.RoleMentor {
		filter.MentorID = &userID
	}
	if req.Status != "" {
		filter.Status = &req.Status
	}

	items, total, err := h.s.ListMentorRequests(c.Request.Context(), filter)
	if err != nil {
		h.mapError(c, err)
		return
	}

	result := make([]dto.MentorRequestResponse, 0, len(items))
	for _, item := range items {
		copy := item
		result = append(result, mapper.ToMentorRequestResponse(&copy))
	}

	c.JSON(http.StatusOK, dto.ListResponse[dto.MentorRequestResponse]{Items: result, TotalCount: total, Limit: req.Limit, Offset: req.Offset})
}

func (h *Handler) ListMentorSlots(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	slots, total, err := h.s.ListSlots(c.Request.Context(), postgres.SlotFilter{MentorID: &mentorID, Limit: 100, Offset: 0})
	if err != nil {
		h.mapError(c, err)
		return
	}

	items := make([]dto.SlotResponse, 0, len(slots))
	for _, slot := range slots {
		copy := slot
		items = append(items, mapper.ToSlotResponse(&copy))
	}
	c.JSON(http.StatusOK, dto.ListResponse[dto.SlotResponse]{Items: items, TotalCount: total, Limit: 100, Offset: 0})
}

func (h *Handler) CreateMentorSlot(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	var req dto.CreateSlotRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	cityID, parseErr := uuid.Parse(req.CityID)
	if parseErr != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	roomID, parseErr := mapper.ParseUUIDPtr(req.RoomID)
	if parseErr != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	status := req.Status
	if status == "" {
		status = domain.SlotStatusActive
	}

	created, err := h.s.CreateSlot(c.Request.Context(), mentorID, mentorID, roomID, cityID, req.Type, req.StartAt, req.EndAt, req.MeetingURL, status, req.Capacity)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.ToSlotResponse(created))
}

func (h *Handler) UpdateMentorSlot(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var req dto.UpdateSlotRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	patch := map[string]any{}
	if req.StartAt != nil {
		patch["start_at"] = *req.StartAt
	}
	if req.EndAt != nil {
		patch["end_at"] = *req.EndAt
	}
	if req.MeetingURL != nil {
		patch["meeting_url"] = *req.MeetingURL
	}
	if req.Status != nil {
		patch["status"] = *req.Status
	}
	if req.Capacity != nil {
		patch["capacity"] = *req.Capacity
	}
	if req.RoomID != nil {
		roomID, parseErr := uuid.Parse(*req.RoomID)
		if parseErr != nil {
			h.mapError(c, domain.ErrInvalidInput)
			return
		}
		patch["room_id"] = roomID
	}

	updated, err := h.s.UpdateSlot(c.Request.Context(), slotID, mentorID, patch)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToSlotResponse(updated))
}

func (h *Handler) DeleteMentorSlot(c *gin.Context) {
	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	if err = h.s.DeleteSlot(c.Request.Context(), slotID); err != nil {
		h.mapError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) ApproveMentorRequest(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	updated, err := h.s.UpdateMentorRequestStatus(c.Request.Context(), requestID, domain.MentorRequestStatusApproved, &mentorID)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToMentorRequestResponse(updated))
}

func (h *Handler) RejectMentorRequest(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	updated, err := h.s.UpdateMentorRequestStatus(c.Request.Context(), requestID, domain.MentorRequestStatusRejected, &mentorID)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToMentorRequestResponse(updated))
}

func (h *Handler) SubscribeMentorSkill(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	var req struct {
		SkillID string `binding:"required,uuid" json:"skill_id"`
	}
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	skillID, parseErr := uuid.Parse(req.SkillID)
	if parseErr != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	if err = h.s.SubscribeMentorSkill(c.Request.Context(), mentorID, skillID); err != nil {
		h.mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UnsubscribeMentorSkill(c *gin.Context) {
	mentorID, err := firstUserID(c)
	if err != nil {
		h.mapError(c, err)
		return
	}
	skillID, err := uuid.Parse(c.Param("skillId"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	if err = h.s.UnsubscribeMentorSkill(c.Request.Context(), mentorID, skillID); err != nil {
		h.mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func timePtr(v *time.Time) *time.Time {
	return v
}

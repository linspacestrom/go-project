package platform

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/dto"
	"github.com/linspacestrom/go-project/internal/mapper"
)

func (h *Handler) RegisterMentor(c *gin.Context) {
	var req dto.RegisterMentorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	cityID, err := uuid.Parse(req.CityID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	user, err := h.s.RegisterMentor(c.Request.Context(), req.Email, req.Password, req.FullName, cityID, req.Description, req.Title)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "email already exists"}})
			return
		}
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.ToUserResponse(user))
}

func (h *Handler) CreateCity(c *gin.Context) {
	var req dto.CreateCityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	city, err := h.s.CreateCity(c.Request.Context(), req.Name, isActive)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.ToCityResponse(city))
}

func (h *Handler) UpdateCity(c *gin.Context) {
	cityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var req dto.UpdateCityRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	city, err := h.s.UpdateCity(c.Request.Context(), cityID, req.Name, req.IsActive)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToCityResponse(city))
}

func (h *Handler) CreateHub(c *gin.Context) {
	var req dto.CreateHubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	cityID, err := uuid.Parse(req.CityID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	hub, err := h.s.CreateHub(c.Request.Context(), cityID, req.Name, req.Address, isActive)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusCreated, mapper.ToHubResponse(hub))
}

func (h *Handler) UpdateHub(c *gin.Context) {
	hubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var req dto.UpdateHubRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	hub, err := h.s.UpdateHub(c.Request.Context(), hubID, req.Name, req.Address, req.IsActive)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToHubResponse(hub))
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	hubID, err := uuid.Parse(req.HubID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	room, err := h.s.CreateRoom(c.Request.Context(), hubID, req.Name, req.Description, req.RoomType, req.Capacity, isActive)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusCreated, mapper.ToRoomResponse(room))
}

func (h *Handler) UpdateRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var req dto.UpdateRoomRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}

	room, err := h.s.UpdateRoom(c.Request.Context(), roomID, req.Name, req.Description, req.RoomType, req.Capacity, req.IsActive)
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapper.ToRoomResponse(room))
}

func (h *Handler) DeleteRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	if err = h.s.DeleteRoom(c.Request.Context(), roomID); err != nil {
		h.mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateUserCity(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	var req dto.UpdateUserCityRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	cityID, err := uuid.Parse(req.CityID)
	if err != nil {
		h.mapError(c, domain.ErrInvalidInput)
		return
	}
	if err = h.s.UpdateUserCity(c.Request.Context(), userID, cityID); err != nil {
		h.mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) BusinessAnalytics(c *gin.Context) {
	stats, err := h.s.GetBusinessAnalytics(c.Request.Context())
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.BusinessAnalyticsResponse{
		MentorApprovals:                stats.MentorApprovals,
		MentorRejections:               stats.MentorRejections,
		ActiveBookings:                 stats.ActiveBookings,
		MentorRequestsBySkill:          stats.MentorRequestsBySkill,
		MentorRequestsByCategory:       stats.MentorRequestsByCategory,
		UserCancelsAfterMentorApproval: stats.UserCancelsAfterMentorApproval,
	})
}

func (h *Handler) TechnicalAnalytics(c *gin.Context) {
	stats, err := h.s.GetTechnicalAnalytics(c.Request.Context())
	if err != nil {
		h.mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.TechnicalAnalyticsResponse{
		BookingConflictCount: stats.BookingConflictCount,
		FailedOutboxEvents:   stats.FailedOutboxEvents,
		OutboxLagCount:       stats.OutboxLagCount,
	})
}

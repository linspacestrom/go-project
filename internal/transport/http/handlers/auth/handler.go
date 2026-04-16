package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/linspacestrom/go-project/internal/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/dto"
	"github.com/linspacestrom/go-project/internal/mapper"
	authService "github.com/linspacestrom/go-project/internal/service/auth"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(ctx context.Context, params authService.RegisterParams) (*domain.User, string, string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, token string) (string, string, error)
	Logout(ctx context.Context, token string) error
}

type Handler struct {
	s AuthService
}

func NewHandler(s AuthService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	auth := r.Group("/api/v1/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/logout", h.Logout)
}

func (h *Handler) Register(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()}})
		return
	}

	course := 0
	if req.Course != nil {
		course = *req.Course
	}
	created, access, refresh, err := h.s.Register(c.Request.Context(), authService.RegisterParams{
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
		FullName:   req.FullName,
		BirthDate:  req.BirthDate,
		University: req.University,
		Course:     course,
		DegreeType: req.DegreeType,
	})
	if err != nil {
		h.mapError(c, log, err, "failed to register")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":   mapper.ToUserResponse(created),
		"tokens": dto.TokenPairResponse{AccessToken: access, RefreshToken: refresh},
	})
}

func (h *Handler) Login(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()}})
		return
	}

	access, refresh, err := h.s.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.mapError(c, log, err, "failed to login")
		return
	}

	c.JSON(http.StatusOK, dto.TokenPairResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *Handler) Refresh(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()}})
		return
	}

	access, refresh, err := h.s.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.mapError(c, log, err, "failed to refresh tokens")
		return
	}

	c.JSON(http.StatusOK, dto.TokenPairResponse{AccessToken: access, RefreshToken: refresh})
}

func (h *Handler) Logout(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()}})
		return
	}

	if err := h.s.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		h.mapError(c, log, err, "failed to logout")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) mapError(c *gin.Context, log *zap.Logger, err error, msg string) {
	switch {
	case errors.Is(err, domain.ErrEmailExists):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "email already exists"}})
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: "invalid credentials"}})
	case errors.Is(err, domain.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "invalid input"}})
	default:
		log.Error(msg, zap.Error(err), zap.String("request_id", uuid.NewString()))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"}})
	}
}

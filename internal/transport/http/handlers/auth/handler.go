package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/linspacestrom/go-project/internal/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/go-project/internal/domain"
	consts "github.com/linspacestrom/go-project/internal/domain"
	"github.com/linspacestrom/go-project/internal/dto"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(ctx context.Context, email, password, role string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type Handler struct {
	s AuthService
}

func NewHandler(s AuthService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.POST("/register/user", h.RegisterUser)
	r.POST("/register/mentor", h.RegisterMentor)
	r.POST("/login", h.Login)
}

func (h *Handler) RegisterUser(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})

		return
	}

	user, err := h.s.Register(c.Request.Context(), req.Email, req.Password, consts.RoleUser)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "email already exists"},
			})

			return
		}

		log.Error("failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (h *Handler) Login(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})

		return
	}

	token, err := h.s.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "UNAUTHORIZED", Message: "invalid credentials"},
			})

			return
		}

		log.Error("failed to login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal error"},
		})

		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: token})
}

func (h *Handler) RegisterMentor(c *gin.Context) {
	log := middleware.GetLoggerFromContext(c)

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: err.Error()},
		})

		return
	}

	user, err := h.s.Register(c.Request.Context(), req.Email, req.Password, consts.RoleMentor)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "email already exists"},
			})

			return
		}

		log.Error("failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "internal server error"},
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

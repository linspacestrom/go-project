package info

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	db Pinger
}

func NewHandler(db Pinger) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/heathz", h.Healthz)
	r.GET("/readyz", h.Readyz)
}

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) Readyz(c *gin.Context) {
	dbStatus := "connected"
	if err := h.db.Ping(c.Request.Context()); err != nil {
		dbStatus = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": dbStatus,
	})
}

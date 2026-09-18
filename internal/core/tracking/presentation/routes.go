package presentation

import "github.com/gin-gonic/gin"

// RegisterRoutes expone los endpoints del modulo tracking.
func RegisterRoutes(router *gin.Engine, h *Handler) {
	group := router.Group("/api/tracking")
	{
		group.POST("/positions", h.Ingest)
		group.GET("/devices/:deviceID/latest", h.Latest)
		group.GET("/devices/:deviceID/history", h.History)
	}
}

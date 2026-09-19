package presentation

import "github.com/gin-gonic/gin"

// RegisterRoutes expone los endpoints del modulo tracking.
// POST /positions va con auth de dispositivo; las consultas con auth de admin.
func RegisterRoutes(router *gin.Engine, h *Handler, adminAuth, deviceAuth gin.HandlerFunc) {
	group := router.Group("/api/tracking")
	{
		group.POST("/positions", deviceAuth, h.Ingest)
		group.GET("/devices/:deviceID/latest", adminAuth, h.Latest)
		group.GET("/devices/:deviceID/history", adminAuth, h.History)
	}
}

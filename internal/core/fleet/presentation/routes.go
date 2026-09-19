package presentation

import "github.com/gin-gonic/gin"

// RegisterRoutes expone los endpoints del modulo fleet.
func RegisterRoutes(router *gin.Engine, h *Handler) {
	group := router.Group("/api/fleet")
	{
		group.POST("/devices", h.CreateDevice)
		group.GET("/devices", h.ListDevices)
		group.GET("/devices/:id", h.GetDevice)

		group.POST("/assets", h.CreateAsset)
		group.GET("/assets", h.ListAssets)
		group.GET("/assets/:id", h.GetAsset)
		group.POST("/assets/:id/assign-device", h.AssignDevice)
		group.POST("/assets/:id/unassign-device", h.UnassignDevice)
	}
}

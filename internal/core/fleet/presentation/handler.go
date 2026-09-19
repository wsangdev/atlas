package presentation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"atlas/internal/core/fleet/application"
	"atlas/internal/core/fleet/domain"
)

type Handler struct {
	createDevice   *application.CreateDevice
	createAsset    *application.CreateAsset
	assignDevice   *application.AssignDevice
	unassignDevice *application.UnassignDevice
	getDevice      *application.GetDevice
	listDevices    *application.ListDevices
	getAsset       *application.GetAsset
	listAssets     *application.ListAssets
}

func NewHandler(
	createDevice *application.CreateDevice,
	createAsset *application.CreateAsset,
	assignDevice *application.AssignDevice,
	unassignDevice *application.UnassignDevice,
	getDevice *application.GetDevice,
	listDevices *application.ListDevices,
	getAsset *application.GetAsset,
	listAssets *application.ListAssets,
) *Handler {
	return &Handler{
		createDevice:   createDevice,
		createAsset:    createAsset,
		assignDevice:   assignDevice,
		unassignDevice: unassignDevice,
		getDevice:      getDevice,
		listDevices:    listDevices,
		getAsset:       getAsset,
		listAssets:     listAssets,
	}
}

type createDeviceRequest struct {
	Serial   string `json:"serial" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Protocol string `json:"protocol"`
}

func (h *Handler) CreateDevice(c *gin.Context) {
	var req createDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload invalido", "error": err.Error()})
		return
	}

	device, err := h.createDevice.Execute(application.CreateDeviceInput{
		Serial:   req.Serial,
		Name:     req.Name,
		Protocol: req.Protocol,
	})
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": deviceResponse(device)})
}

func (h *Handler) ListDevices(c *gin.Context) {
	devices, err := h.listDevices.Execute()
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	data := make([]gin.H, 0, len(devices))
	for _, d := range devices {
		data = append(data, deviceResponse(d))
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) GetDevice(c *gin.Context) {
	device, err := h.getDevice.Execute(c.Param("id"))
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": deviceResponse(*device)})
}

type createAssetRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Plate string `json:"plate"`
}

func (h *Handler) CreateAsset(c *gin.Context) {
	var req createAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload invalido", "error": err.Error()})
		return
	}

	asset, err := h.createAsset.Execute(application.CreateAssetInput{
		Name:  req.Name,
		Type:  req.Type,
		Plate: req.Plate,
	})
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": assetResponse(asset)})
}

func (h *Handler) ListAssets(c *gin.Context) {
	assets, err := h.listAssets.Execute()
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	data := make([]gin.H, 0, len(assets))
	for _, a := range assets {
		data = append(data, assetResponse(a))
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *Handler) GetAsset(c *gin.Context) {
	asset, err := h.getAsset.Execute(c.Param("id"))
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assetResponse(*asset)})
}

type assignDeviceRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

func (h *Handler) AssignDevice(c *gin.Context) {
	var req assignDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload invalido", "error": err.Error()})
		return
	}

	asset, err := h.assignDevice.Execute(application.AssignDeviceInput{
		AssetID:  c.Param("id"),
		DeviceID: req.DeviceID,
	})
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assetResponse(asset)})
}

func (h *Handler) UnassignDevice(c *gin.Context) {
	asset, err := h.unassignDevice.Execute(c.Param("id"))
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": assetResponse(asset)})
}

func deviceResponse(d domain.Device) gin.H {
	return gin.H{
		"id":         d.ID,
		"serial":     d.Serial,
		"name":       d.Name,
		"protocol":   d.Protocol,
		"active":     d.Active,
		"created_at": d.CreatedAt,
		"updated_at": d.UpdatedAt,
	}
}

func assetResponse(a domain.Asset) gin.H {
	return gin.H{
		"id":         a.ID,
		"name":       a.Name,
		"type":       string(a.Type),
		"plate":      a.Plate,
		"device_id":  a.DeviceID,
		"created_at": a.CreatedAt,
		"updated_at": a.UpdatedAt,
	}
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrDeviceNotFound), errors.Is(err, domain.ErrAssetNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrSerialDuplicated),
		errors.Is(err, domain.ErrDeviceAlreadyAssigned),
		errors.Is(err, domain.ErrDeviceInactive),
		errors.Is(err, domain.ErrDeviceNotAssigned):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrSerialRequired),
		errors.Is(err, domain.ErrDeviceNameRequired),
		errors.Is(err, domain.ErrAssetNameRequired),
		errors.Is(err, domain.ErrInvalidAssetType):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "error interno"
	}
}

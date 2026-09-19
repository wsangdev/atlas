package presentation

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"atlas/internal/core/tracking/application"
	"atlas/internal/core/tracking/domain"
)

type Handler struct {
	ingest  *application.IngestPosition
	latest  *application.GetLatestPosition
	history *application.GetHistory
}

func NewHandler(
	ingest *application.IngestPosition,
	latest *application.GetLatestPosition,
	history *application.GetHistory,
) *Handler {
	return &Handler{ingest: ingest, latest: latest, history: history}
}

type ingestPositionRequest struct {
	DeviceID   string     `json:"device_id" binding:"required"`
	Lat        float64    `json:"lat"`
	Lng        float64    `json:"lng"`
	SpeedKmh   *float64   `json:"speed_kmh"`
	Heading    *float64   `json:"heading"`
	AccuracyM  *float64   `json:"accuracy_m"`
	Source     string     `json:"source"`
	RecordedAt *time.Time `json:"recorded_at"`
}

// Ingest recibe una posicion nueva (POST /api/tracking/positions).
func (h *Handler) Ingest(c *gin.Context) {
	var req ingestPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "payload invalido", "error": err.Error()})
		return
	}

	position, err := h.ingest.Execute(application.IngestPositionInput{
		DeviceID:   req.DeviceID,
		Lat:        req.Lat,
		Lng:        req.Lng,
		SpeedKmh:   req.SpeedKmh,
		Heading:    req.Heading,
		AccuracyM:  req.AccuracyM,
		Source:     req.Source,
		RecordedAt: req.RecordedAt,
	})
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": positionResponse(position)})
}

// Latest devuelve la ultima posicion de un dispositivo.
func (h *Handler) Latest(c *gin.Context) {
	position, err := h.latest.Execute(c.Param("deviceID"))
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}
	if position == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "sin posiciones para el dispositivo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": positionResponse(*position)})
}

// History devuelve el recorrido de un dispositivo en un rango.
func (h *Handler) History(c *gin.Context) {
	now := time.Now().UTC()
	from := now.Add(-24 * time.Hour)
	to := now

	if raw := c.Query("from"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "from invalido, usa RFC3339"})
			return
		}
		from = parsed
	}
	if raw := c.Query("to"); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "to invalido, usa RFC3339"})
			return
		}
		to = parsed
	}

	limit, _ := strconv.Atoi(c.Query("limit"))

	positions, err := h.history.Execute(application.GetHistoryInput{
		DeviceID: c.Param("deviceID"),
		From:     from,
		To:       to,
		Limit:    limit,
	})
	if err != nil {
		status, message := mapError(err)
		c.JSON(status, gin.H{"message": message})
		return
	}

	data := make([]gin.H, 0, len(positions))
	for _, p := range positions {
		data = append(data, positionResponse(p))
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func positionResponse(p domain.Position) gin.H {
	return gin.H{
		"id":          p.ID,
		"device_id":   p.DeviceID,
		"lat":         p.Coords.Lat,
		"lng":         p.Coords.Long,
		"speed_kmh":   p.SpeedKmh,
		"heading":     p.Heading,
		"accuracy_m":  p.AccuracyM,
		"source":      p.Source,
		"recorded_at": p.RecordedAt,
		"created_at":  p.CreatedAt,
	}
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrDeviceNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrDeviceInactive):
		return http.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrDeviceRequired),
		errors.Is(err, domain.ErrInvalidDeviceID),
		errors.Is(err, domain.ErrInvalidCoordinates),
		errors.Is(err, application.ErrInvalidRange):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, "error interno"
	}
}

package infrastructure

import (
	"time"

	"atlas/internal/core/tracking/domain"
	"atlas/internal/shared/ws"
)

// positionPayload es lo que viaja por WebSocket (plano, sin anidar Coords).
type positionPayload struct {
	ID         string    `json:"id"`
	DeviceID   string    `json:"device_id"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	SpeedKmh   *float64  `json:"speed_kmh,omitempty"`
	Heading    *float64  `json:"heading,omitempty"`
	AccuracyM  *float64  `json:"accuracy_m,omitempty"`
	Source     string    `json:"source"`
	RecordedAt time.Time `json:"recorded_at"`
}

// WSPublisher implementa application.PositionPublisher sobre el hub compartido.
type WSPublisher struct {
	hub *ws.Hub
}

func NewWSPublisher(hub *ws.Hub) *WSPublisher {
	return &WSPublisher{hub: hub}
}

func (p *WSPublisher) PublishPosition(position domain.Position) {
	p.hub.Broadcast("position.created", positionPayload{
		ID:         position.ID,
		DeviceID:   position.DeviceID,
		Lat:        position.Coords.Lat,
		Lng:        position.Coords.Long,
		SpeedKmh:   position.SpeedKmh,
		Heading:    position.Heading,
		AccuracyM:  position.AccuracyM,
		Source:     position.Source,
		RecordedAt: position.RecordedAt,
	})
}

package infrastructure

import (
	"time"

	"atlas/internal/core/tracking/domain"
)

// PositionModel es la representacion de la tabla "positions" en Postgres.
// Es el unico lugar donde viven las etiquetas de GORM.
type PositionModel struct {
	ID         string    `gorm:"type:uuid;primaryKey"`
	DeviceID   string    `gorm:"size:64;not null;index:idx_positions_device_recorded,priority:1"`
	Lat        float64   `gorm:"not null"`
	Lng        float64   `gorm:"not null"`
	SpeedKmh   *float64  `gorm:"column:speed_kmh"`
	Heading    *float64  `gorm:"column:heading"`
	AccuracyM  *float64  `gorm:"column:accuracy_m"`
	Source     string    `gorm:"size:40;index"`
	RecordedAt time.Time `gorm:"not null;index:idx_positions_device_recorded,priority:2"`
	CreatedAt  time.Time `gorm:"not null"`
}

func (PositionModel) TableName() string { return "positions" }

func toModel(p domain.Position) PositionModel {
	return PositionModel{
		ID:         p.ID,
		DeviceID:   p.DeviceID,
		Lat:        p.Coords.Lat,
		Lng:        p.Coords.Long,
		SpeedKmh:   p.SpeedKmh,
		Heading:    p.Heading,
		AccuracyM:  p.AccuracyM,
		Source:     p.Source,
		RecordedAt: p.RecordedAt,
		CreatedAt:  p.CreatedAt,
	}
}

func toDomain(m PositionModel) domain.Position {
	return domain.Position{
		ID:         m.ID,
		DeviceID:   m.DeviceID,
		Coords:     domain.Coordinates{Lat: m.Lat, Long: m.Lng},
		SpeedKmh:   m.SpeedKmh,
		Heading:    m.Heading,
		AccuracyM:  m.AccuracyM,
		Source:     m.Source,
		RecordedAt: m.RecordedAt,
		CreatedAt:  m.CreatedAt,
	}
}

package domain

import (
	"errors"
	"time"
)

var (
	ErrDeviceRequired     = errors.New("device_id es requerido")
	ErrInvalidCoordinates = errors.New("coordenadas invalidas")
)

type Coordinates struct {
	Lat  float64
	Long float64
}

func NewCoordinates(lat, long float64) (Coordinates, error) {
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return Coordinates{}, ErrInvalidCoordinates
	}

	if lat == 0 && long == 0 {
		return Coordinates{}, ErrInvalidCoordinates
	}

	return Coordinates{Lat: lat, Long: long}, nil
}

// Position es la entidad central del modulo tracking.
type Position struct {
	ID         string
	DeviceID   string
	Coords     Coordinates
	SpeedKmh   *float64
	Heading    *float64
	AccuracyM  *float64
	Source     string
	RecordedAt time.Time
	CreatedAt  time.Time
}

type PositionRepository interface {
	Save(p Position) error
	LatestByDevice(deviceID string) (*Position, error)
	ListByDeviceRange(deviceID string, from, to time.Time, limit int) ([]Position, error)
}

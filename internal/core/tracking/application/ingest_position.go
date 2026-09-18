package application

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"atlas/internal/core/tracking/domain"
)

// IngestPositionInput es lo que recibe el caso de uso desde la capa de
// presentacion (HTTP, MQTT, TCP...). No depende de Gin ni de JSON.
type IngestPositionInput struct {
	DeviceID   string
	Lat        float64
	Lng        float64
	SpeedKmh   *float64
	Heading    *float64
	AccuracyM  *float64
	Source     string
	RecordedAt *time.Time
}

// IngestPosition es el caso de uso central: valida, guarda y publica.
type IngestPosition struct {
	repo      domain.PositionRepository
	publisher PositionPublisher
	now       func() time.Time
}

func NewIngestPosition(repo domain.PositionRepository, publisher PositionPublisher) *IngestPosition {
	return &IngestPosition{
		repo:      repo,
		publisher: publisher,
		now:       time.Now,
	}
}

func (uc *IngestPosition) Execute(input IngestPositionInput) (domain.Position, error) {
	deviceID := strings.TrimSpace(input.DeviceID)
	if deviceID == "" {
		return domain.Position{}, domain.ErrDeviceRequired
	}

	coords, err := domain.NewCoordinates(input.Lat, input.Lng)
	if err != nil {
		return domain.Position{}, err
	}

	now := uc.now()
	recordedAt := now
	if input.RecordedAt != nil && !input.RecordedAt.IsZero() {
		recordedAt = input.RecordedAt.UTC()
	}

	source := strings.TrimSpace(input.Source)
	if source == "" {
		source = "http"
	}

	id, err := uuid.NewV7()
	if err != nil {
		return domain.Position{}, err
	}

	position := domain.Position{
		ID:         id.String(),
		DeviceID:   deviceID,
		Coords:     coords,
		SpeedKmh:   input.SpeedKmh,
		Heading:    input.Heading,
		AccuracyM:  input.AccuracyM,
		Source:     source,
		RecordedAt: recordedAt,
		CreatedAt:  now,
	}

	if err := uc.repo.Save(position); err != nil {
		return domain.Position{}, err
	}

	// Best-effort: si el aviso en vivo falla, la posicion ya quedo guardada.
	uc.publisher.PublishPosition(position)

	return position, nil
}

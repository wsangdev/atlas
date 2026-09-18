package application

import (
	"strings"
	"time"

	"atlas/internal/core/tracking/domain"
)

const (
	defaultHistoryLimit = 500
	maxHistoryLimit     = 5000
)

type GetHistoryInput struct {
	DeviceID string
	From     time.Time
	To       time.Time
	Limit    int
}

type GetHistory struct {
	repo domain.PositionRepository
}

func NewGetHistory(repo domain.PositionRepository) *GetHistory {
	return &GetHistory{repo: repo}
}

func (uc *GetHistory) Execute(input GetHistoryInput) ([]domain.Position, error) {
	deviceID := strings.TrimSpace(input.DeviceID)
	if deviceID == "" {
		return nil, domain.ErrDeviceRequired
	}
	if input.From.IsZero() || input.To.IsZero() || !input.From.Before(input.To) {
		return nil, ErrInvalidRange
	}

	limit := input.Limit
	if limit <= 0 {
		limit = defaultHistoryLimit
	}
	if limit > maxHistoryLimit {
		limit = maxHistoryLimit
	}

	return uc.repo.ListByDeviceRange(deviceID, input.From.UTC(), input.To.UTC(), limit)
}

package application

import (
	"strings"

	"github.com/gofrs/uuid/v5"

	"atlas/internal/core/tracking/domain"
)

type GetLatestPosition struct {
	repo domain.PositionRepository
}

func NewGetLatestPosition(repo domain.PositionRepository) *GetLatestPosition {
	return &GetLatestPosition{repo: repo}
}

func (uc *GetLatestPosition) Execute(deviceID string) (*domain.Position, error) {
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return nil, domain.ErrDeviceRequired
	}
	if _, err := uuid.FromString(deviceID); err != nil {
		return nil, domain.ErrInvalidDeviceID
	}
	return uc.repo.LatestByDevice(deviceID)
}

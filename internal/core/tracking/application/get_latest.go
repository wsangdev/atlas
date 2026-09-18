package application

import (
	"strings"

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
	return uc.repo.LatestByDevice(deviceID)
}

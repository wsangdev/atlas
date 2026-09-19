package application

import (
	"strings"
	"time"

	"atlas/internal/core/fleet/domain"
)

type UnassignDevice struct {
	assets domain.AssetRepository
	now    func() time.Time
}

func NewUnassignDevice(assets domain.AssetRepository) *UnassignDevice {
	return &UnassignDevice{assets: assets, now: time.Now}
}

func (uc *UnassignDevice) Execute(assetID string) (domain.Asset, error) {
	asset, err := uc.assets.FindByID(strings.TrimSpace(assetID))
	if err != nil {
		return domain.Asset{}, err
	}
	if asset == nil {
		return domain.Asset{}, domain.ErrAssetNotFound
	}
	if asset.DeviceID == nil {
		return domain.Asset{}, domain.ErrDeviceNotAssigned
	}

	asset.DeviceID = nil
	asset.UpdatedAt = uc.now()

	if err := uc.assets.Update(*asset); err != nil {
		return domain.Asset{}, err
	}
	return *asset, nil
}

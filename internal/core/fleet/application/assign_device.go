package application

import (
	"strings"
	"time"

	"atlas/internal/core/fleet/domain"
)

type AssignDeviceInput struct {
	AssetID  string
	DeviceID string
}

// AssignDevice asigna un dispositivo a un activo. Regla de negocio: un
// dispositivo solo puede estar en un activo a la vez.
type AssignDevice struct {
	assets  domain.AssetRepository
	devices domain.DeviceRepository
	now     func() time.Time
}

func NewAssignDevice(assets domain.AssetRepository, devices domain.DeviceRepository) *AssignDevice {
	return &AssignDevice{assets: assets, devices: devices, now: time.Now}
}

func (uc *AssignDevice) Execute(input AssignDeviceInput) (domain.Asset, error) {
	asset, err := uc.assets.FindByID(strings.TrimSpace(input.AssetID))
	if err != nil {
		return domain.Asset{}, err
	}
	if asset == nil {
		return domain.Asset{}, domain.ErrAssetNotFound
	}

	device, err := uc.devices.FindByID(strings.TrimSpace(input.DeviceID))
	if err != nil {
		return domain.Asset{}, err
	}
	if device == nil {
		return domain.Asset{}, domain.ErrDeviceNotFound
	}
	if !device.Active {
		return domain.Asset{}, domain.ErrDeviceInactive
	}

	assignedAsset, err := uc.assets.FindByDeviceID(device.ID)
	if err != nil {
		return domain.Asset{}, err
	}
	if assignedAsset != nil && assignedAsset.ID != asset.ID {
		return domain.Asset{}, domain.ErrDeviceAlreadyAssigned
	}

	asset.DeviceID = &device.ID
	asset.UpdatedAt = uc.now()

	if err := uc.assets.Update(*asset); err != nil {
		return domain.Asset{}, err
	}
	return *asset, nil
}

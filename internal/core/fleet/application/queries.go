package application

import (
	"strings"

	"atlas/internal/core/fleet/domain"
)

// ── Consultas (queries) ──

type GetDevice struct {
	repo domain.DeviceRepository
}

func NewGetDevice(repo domain.DeviceRepository) *GetDevice {
	return &GetDevice{repo: repo}
}

func (q *GetDevice) Execute(id string) (*domain.Device, error) {
	device, err := q.repo.FindByID(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, domain.ErrDeviceNotFound
	}
	return device, nil
}

type ListDevices struct {
	repo domain.DeviceRepository
}

func NewListDevices(repo domain.DeviceRepository) *ListDevices {
	return &ListDevices{repo: repo}
}

func (q *ListDevices) Execute() ([]domain.Device, error) {
	return q.repo.List()
}

type GetAsset struct {
	repo domain.AssetRepository
}

func NewGetAsset(repo domain.AssetRepository) *GetAsset {
	return &GetAsset{repo: repo}
}

func (q *GetAsset) Execute(id string) (*domain.Asset, error) {
	asset, err := q.repo.FindByID(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, domain.ErrAssetNotFound
	}
	return asset, nil
}

type ListAssets struct {
	repo domain.AssetRepository
}

func NewListAssets(repo domain.AssetRepository) *ListAssets {
	return &ListAssets{repo: repo}
}

func (q *ListAssets) Execute() ([]domain.Asset, error) {
	return q.repo.List()
}

package domain

import (
	"strings"
	"time"
)

// AssetType es el tipo de activo rastreado.
type AssetType string

const (
	AssetTypeMoto   AssetType = "moto"
	AssetTypeAuto   AssetType = "auto"
	AssetTypeObjeto AssetType = "objeto"
)

func ParseAssetType(raw string) (AssetType, error) {
	switch AssetType(strings.ToLower(strings.TrimSpace(raw))) {
	case AssetTypeMoto:
		return AssetTypeMoto, nil
	case AssetTypeAuto:
		return AssetTypeAuto, nil
	case AssetTypeObjeto:
		return AssetTypeObjeto, nil
	default:
		return "", ErrInvalidAssetType
	}
}

// Asset es lo que se rastrea: una moto, un auto u otro objeto. Un activo tiene
// a lo sumo UN dispositivo asignado (relacion 1:1 logica).
type Asset struct {
	ID        string
	Name      string
	Type      AssetType
	Plate     string
	DeviceID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AssetRepository interface {
	Save(a Asset) error
	Update(a Asset) error
	FindByID(id string) (*Asset, error)
	FindByDeviceID(deviceID string) (*Asset, error)
	List() ([]Asset, error)
}

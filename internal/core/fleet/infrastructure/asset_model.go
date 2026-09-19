package infrastructure

import (
	"time"

	"atlas/internal/core/fleet/domain"
)

type AssetModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:120;not null"`
	Type      string    `gorm:"size:20;not null;index"`
	Plate     string    `gorm:"size:20;index"`
	DeviceID  *string   `gorm:"type:uuid;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (AssetModel) TableName() string { return "assets" }

func assetToModel(a domain.Asset) AssetModel {
	return AssetModel{
		ID:        a.ID,
		Name:      a.Name,
		Type:      string(a.Type),
		Plate:     a.Plate,
		DeviceID:  a.DeviceID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func assetToDomain(m AssetModel) domain.Asset {
	return domain.Asset{
		ID:        m.ID,
		Name:      m.Name,
		Type:      domain.AssetType(m.Type),
		Plate:     m.Plate,
		DeviceID:  m.DeviceID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

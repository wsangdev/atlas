package infrastructure

import (
	"time"

	"atlas/internal/core/fleet/domain"
)

type DeviceModel struct {
	ID         string    `gorm:"type:uuid;primaryKey"`
	Serial     string    `gorm:"size:80;not null;uniqueIndex"`
	Name       string    `gorm:"size:120;not null"`
	Protocol   string    `gorm:"size:40;not null"`
	Active     bool      `gorm:"not null;default:true"`
	APIKeyHash string    `gorm:"size:64"` // nullable: devices previos no tienen key
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

func (DeviceModel) TableName() string { return "devices" }

func deviceToModel(d domain.Device) DeviceModel {
	return DeviceModel{
		ID:         d.ID,
		Serial:     d.Serial,
		Name:       d.Name,
		Protocol:   d.Protocol,
		Active:     d.Active,
		APIKeyHash: d.APIKeyHash,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

func deviceToDomain(m DeviceModel) domain.Device {
	return domain.Device{
		ID:         m.ID,
		Serial:     m.Serial,
		Name:       m.Name,
		Protocol:   m.Protocol,
		Active:     m.Active,
		APIKeyHash: m.APIKeyHash,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

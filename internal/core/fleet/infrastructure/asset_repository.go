package infrastructure

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"atlas/internal/core/fleet/domain"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Save(a domain.Asset) error {
	model := assetToModel(a)
	if err := r.db.Create(&model).Error; err != nil {
		return fmt.Errorf("guardar activo: %w", err)
	}
	return nil
}

func (r *AssetRepository) Update(a domain.Asset) error {
	model := assetToModel(a)
	if err := r.db.Model(&AssetModel{}).Where("id = ?", a.ID).Updates(map[string]any{
		"name":       model.Name,
		"type":       model.Type,
		"plate":      model.Plate,
		"device_id":  model.DeviceID,
		"updated_at": model.UpdatedAt,
	}).Error; err != nil {
		return fmt.Errorf("actualizar activo: %w", err)
	}
	return nil
}

func (r *AssetRepository) FindByID(id string) (*domain.Asset, error) {
	var model AssetModel
	err := r.db.Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar activo: %w", err)
	}
	a := assetToDomain(model)
	return &a, nil
}

func (r *AssetRepository) FindByDeviceID(deviceID string) (*domain.Asset, error) {
	var model AssetModel
	err := r.db.Where("device_id = ?", deviceID).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar activo por dispositivo: %w", err)
	}
	a := assetToDomain(model)
	return &a, nil
}

func (r *AssetRepository) List() ([]domain.Asset, error) {
	var models []AssetModel
	if err := r.db.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listar activos: %w", err)
	}
	out := make([]domain.Asset, 0, len(models))
	for _, m := range models {
		out = append(out, assetToDomain(m))
	}
	return out, nil
}

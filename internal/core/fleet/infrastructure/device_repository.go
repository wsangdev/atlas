package infrastructure

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"atlas/internal/core/fleet/domain"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Save(d domain.Device) error {
	model := deviceToModel(d)
	if err := r.db.Create(&model).Error; err != nil {
		return fmt.Errorf("guardar dispositivo: %w", err)
	}
	return nil
}

func (r *DeviceRepository) Update(d domain.Device) error {
	model := deviceToModel(d)
	if err := r.db.Model(&DeviceModel{}).Where("id = ?", d.ID).Updates(map[string]any{
		"serial":     model.Serial,
		"name":       model.Name,
		"protocol":   model.Protocol,
		"active":     model.Active,
		"updated_at": model.UpdatedAt,
	}).Error; err != nil {
		return fmt.Errorf("actualizar dispositivo: %w", err)
	}
	return nil
}

func (r *DeviceRepository) FindByID(id string) (*domain.Device, error) {
	var model DeviceModel
	err := r.db.Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar dispositivo: %w", err)
	}
	d := deviceToDomain(model)
	return &d, nil
}

func (r *DeviceRepository) FindBySerial(serial string) (*domain.Device, error) {
	var model DeviceModel
	err := r.db.Where("serial = ?", serial).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar dispositivo por serial: %w", err)
	}
	d := deviceToDomain(model)
	return &d, nil
}

func (r *DeviceRepository) List() ([]domain.Device, error) {
	var models []DeviceModel
	if err := r.db.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listar dispositivos: %w", err)
	}
	out := make([]domain.Device, 0, len(models))
	for _, m := range models {
		out = append(out, deviceToDomain(m))
	}
	return out, nil
}

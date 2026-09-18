package infrastructure

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"atlas/internal/core/tracking/domain"
)

// PositionRepository es el adaptador Postgres del puerto domain.PositionRepository.
type PositionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) Save(p domain.Position) error {
	model := toModel(p)
	if err := r.db.Create(&model).Error; err != nil {
		return fmt.Errorf("guardar posicion: %w", err)
	}
	return nil
}

func (r *PositionRepository) LatestByDevice(deviceID string) (*domain.Position, error) {
	var model PositionModel
	err := r.db.
		Where("device_id = ?", deviceID).
		Order("recorded_at DESC").
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("buscar ultima posicion: %w", err)
	}
	p := toDomain(model)
	return &p, nil
}

func (r *PositionRepository) ListByDeviceRange(deviceID string, from, to time.Time, limit int) ([]domain.Position, error) {
	var models []PositionModel
	err := r.db.
		Where("device_id = ? AND recorded_at >= ? AND recorded_at < ?", deviceID, from, to).
		Order("recorded_at ASC").
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("listar posiciones: %w", err)
	}

	out := make([]domain.Position, 0, len(models))
	for _, m := range models {
		out = append(out, toDomain(m))
	}
	return out, nil
}

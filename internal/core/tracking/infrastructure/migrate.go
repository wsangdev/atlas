package infrastructure

import "gorm.io/gorm"

// Migrate crea/actualiza las tablas del modulo tracking.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&PositionModel{})
}

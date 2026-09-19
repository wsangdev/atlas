package infrastructure

import "gorm.io/gorm"

// Migrate crea/actualiza las tablas del modulo fleet.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&DeviceModel{}, &AssetModel{})
}

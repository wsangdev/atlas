package fleet

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"atlas/internal/core/fleet/application"
	"atlas/internal/core/fleet/infrastructure"
	"atlas/internal/core/fleet/presentation"
)

// Module es el composition root del modulo fleet (dispositivos + activos).
type Module struct {
	handler *presentation.Handler
}

func New(db *gorm.DB) (*Module, error) {
	if err := infrastructure.Migrate(db); err != nil {
		return nil, err
	}

	devices := infrastructure.NewDeviceRepository(db)
	assets := infrastructure.NewAssetRepository(db)

	handler := presentation.NewHandler(
		application.NewCreateDevice(devices),
		application.NewCreateAsset(assets),
		application.NewAssignDevice(assets, devices),
		application.NewUnassignDevice(assets),
		application.NewGetDevice(devices),
		application.NewListDevices(devices),
		application.NewGetAsset(assets),
		application.NewListAssets(assets),
	)

	return &Module{handler: handler}, nil
}

func (m *Module) Register(router *gin.Engine) {
	presentation.RegisterRoutes(router, m.handler)
}

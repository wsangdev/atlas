package tracking

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"atlas/internal/core/tracking/application"
	"atlas/internal/core/tracking/infrastructure"
	"atlas/internal/core/tracking/presentation"
	"atlas/internal/shared/ws"
)

// Module es el composition root del modulo tracking: aqui (y solo aqui) se
// conocen todas las capas y se inyectan las dependencias.
type Module struct {
	handler *presentation.Handler
	hub     *ws.Hub
}

func New(db *gorm.DB, hub *ws.Hub, checker application.DeviceChecker) (*Module, error) {
	if err := infrastructure.Migrate(db); err != nil {
		return nil, err
	}

	repo := infrastructure.NewPositionRepository(db)
	publisher := infrastructure.NewWSPublisher(hub)

	ingest := application.NewIngestPosition(repo, publisher, checker)
	latest := application.NewGetLatestPosition(repo)
	history := application.NewGetHistory(repo)

	return &Module{
		handler: presentation.NewHandler(ingest, latest, history),
		hub:     hub,
	}, nil
}

// Register expone las rutas del modulo (HTTP + WebSocket).
func (m *Module) Register(router *gin.Engine, adminAuth, deviceAuth gin.HandlerFunc) {
	presentation.RegisterRoutes(router, m.handler, adminAuth, deviceAuth)

	router.GET("/ws/tracking", func(c *gin.Context) {
		_ = m.hub.Melody().HandleRequest(c.Writer, c.Request)
	})
}

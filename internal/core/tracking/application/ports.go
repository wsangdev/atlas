package application

import "atlas/internal/core/tracking/domain"

// PositionPublisher es un puerto de SALIDA: el caso de uso avisa "llego una
// posicion nueva" sin saber si la implementacion es WebSocket, NATS, etc.
// La implementacion concreta vive en infrastructure.
type PositionPublisher interface {
	PublishPosition(p domain.Position)
}

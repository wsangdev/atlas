package application

import "atlas/internal/core/tracking/domain"

// PositionPublisher es un puerto de SALIDA: el caso de uso avisa "llego una
// posicion nueva" sin saber si la implementacion es WebSocket, NATS, etc.
// La implementacion concreta vive en infrastructure.
type PositionPublisher interface {
	PublishPosition(p domain.Position)
}

// DeviceInfo es el resultado de validar un dispositivo contra el modulo fleet.
type DeviceInfo struct {
	Exists bool
	Active bool
}

// DeviceChecker es un puerto hacia OTRO modulo (fleet). Tracking define que
// necesita; el adaptador que lo implementa se inyecta desde app.go. Asi
// tracking no importa fleet y ambos modulos quedan desacoplados.
type DeviceChecker interface {
	CheckDevice(deviceID string) (DeviceInfo, error)
}

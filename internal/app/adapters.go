package app

import (
	"github.com/gofrs/uuid/v5"

	fleetinfra "atlas/internal/core/fleet/infrastructure"
	"atlas/internal/core/tracking/application"
	"atlas/internal/shared/auth"
)

// trackingDeviceChecker adapta el modulo fleet al puerto DeviceChecker que
// define tracking. Vive en app (composition root) porque es el pegamento entre
// modulos: tracking no importa fleet ni fleet importa tracking.
type trackingDeviceChecker struct {
	devices *fleetinfra.DeviceRepository
}

func (a trackingDeviceChecker) CheckDevice(deviceID string) (application.DeviceInfo, error) {
	device, err := a.devices.FindByID(deviceID)
	if err != nil {
		return application.DeviceInfo{}, err
	}
	if device == nil {
		return application.DeviceInfo{Exists: false}, nil
	}
	return application.DeviceInfo{Exists: true, Active: device.Active}, nil
}

// deviceKeyVerifier implementa auth.DeviceKeyVerifier usando el repo de fleet.
type deviceKeyVerifier struct {
	devices *fleetinfra.DeviceRepository
}

func (v deviceKeyVerifier) VerifyDeviceKey(deviceID, apiKey string) (bool, error) {
	if _, err := uuid.FromString(deviceID); err != nil {
		return false, nil // formato invalido = no autorizado
	}
	device, err := v.devices.FindByID(deviceID)
	if err != nil {
		return false, err
	}
	if device == nil {
		return false, nil
	}
	return auth.VerifyAPIKey(device.APIKeyHash, apiKey), nil
}

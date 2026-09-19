package app

import (
	fleetinfra "atlas/internal/core/fleet/infrastructure"
	"atlas/internal/core/tracking/application"
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

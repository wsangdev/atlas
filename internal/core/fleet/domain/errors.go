package domain

import "errors"

var (
	ErrDeviceNotFound        = errors.New("dispositivo no encontrado")
	ErrAssetNotFound         = errors.New("activo no encontrado")
	ErrInvalidID             = errors.New("id invalido")
	ErrSerialRequired        = errors.New("serial es requerido")
	ErrSerialDuplicated      = errors.New("serial ya registrado")
	ErrDeviceNameRequired    = errors.New("nombre del dispositivo es requerido")
	ErrDeviceInactive        = errors.New("el dispositivo esta inactivo")
	ErrAssetNameRequired     = errors.New("nombre del activo es requerido")
	ErrInvalidAssetType      = errors.New("tipo de activo invalido (moto, auto u objeto)")
	ErrDeviceAlreadyAssigned = errors.New("el dispositivo ya esta asignado a otro activo")
	ErrDeviceNotAssigned     = errors.New("el activo no tiene dispositivo asignado")
)

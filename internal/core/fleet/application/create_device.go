package application

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"atlas/internal/core/fleet/domain"
	"atlas/internal/shared/auth"
)

type CreateDeviceInput struct {
	Serial   string
	Name     string
	Protocol string
}

// CreateDeviceResult devuelve el dispositivo y la API key EN CLARO. Es la
// unica vez que la key viaja: en la DB solo queda su hash.
type CreateDeviceResult struct {
	Device domain.Device
	APIKey string
}

type CreateDevice struct {
	repo domain.DeviceRepository
	now  func() time.Time
}

func NewCreateDevice(repo domain.DeviceRepository) *CreateDevice {
	return &CreateDevice{repo: repo, now: time.Now}
}

func (uc *CreateDevice) Execute(input CreateDeviceInput) (CreateDeviceResult, error) {
	serial := strings.TrimSpace(input.Serial)
	if serial == "" {
		return CreateDeviceResult{}, domain.ErrSerialRequired
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return CreateDeviceResult{}, domain.ErrDeviceNameRequired
	}

	existing, err := uc.repo.FindBySerial(serial)
	if err != nil {
		return CreateDeviceResult{}, err
	}
	if existing != nil {
		return CreateDeviceResult{}, domain.ErrSerialDuplicated
	}

	protocol := strings.TrimSpace(input.Protocol)
	if protocol == "" {
		protocol = "http"
	}

	id, err := uuid.NewV7()
	if err != nil {
		return CreateDeviceResult{}, err
	}

	apiKey, err := auth.GenerateAPIKey()
	if err != nil {
		return CreateDeviceResult{}, err
	}

	now := uc.now()
	device := domain.Device{
		ID:         id.String(),
		Serial:     serial,
		Name:       name,
		Protocol:   protocol,
		Active:     true,
		APIKeyHash: auth.HashAPIKey(apiKey),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.repo.Save(device); err != nil {
		return CreateDeviceResult{}, err
	}
	return CreateDeviceResult{Device: device, APIKey: apiKey}, nil
}

package application

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"atlas/internal/core/fleet/domain"
)

type CreateDeviceInput struct {
	Serial   string
	Name     string
	Protocol string
}

type CreateDevice struct {
	repo domain.DeviceRepository
	now  func() time.Time
}

func NewCreateDevice(repo domain.DeviceRepository) *CreateDevice {
	return &CreateDevice{repo: repo, now: time.Now}
}

func (uc *CreateDevice) Execute(input CreateDeviceInput) (domain.Device, error) {
	serial := strings.TrimSpace(input.Serial)
	if serial == "" {
		return domain.Device{}, domain.ErrSerialRequired
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.Device{}, domain.ErrDeviceNameRequired
	}

	existing, err := uc.repo.FindBySerial(serial)
	if err != nil {
		return domain.Device{}, err
	}
	if existing != nil {
		return domain.Device{}, domain.ErrSerialDuplicated
	}

	protocol := strings.TrimSpace(input.Protocol)
	if protocol == "" {
		protocol = "http"
	}

	id, err := uuid.NewV7()
	if err != nil {
		return domain.Device{}, err
	}

	now := uc.now()
	device := domain.Device{
		ID:        id.String(),
		Serial:    serial,
		Name:      name,
		Protocol:  protocol,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.repo.Save(device); err != nil {
		return domain.Device{}, err
	}
	return device, nil
}

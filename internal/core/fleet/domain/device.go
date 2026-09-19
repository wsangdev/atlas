package domain

import "time"

// Device es el rastreador fisico (IMEI/serial). Hoy la ingesta es por HTTP;
// manana puede ser MQTT o TCP y el Device no cambia.
type Device struct {
	ID        string
	Serial    string
	Name      string
	Protocol  string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeviceRepository interface {
	Save(d Device) error
	Update(d Device) error
	FindByID(id string) (*Device, error)
	FindBySerial(serial string) (*Device, error)
	List() ([]Device, error)
}

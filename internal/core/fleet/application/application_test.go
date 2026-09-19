package application

import (
	"testing"

	"atlas/internal/core/fleet/domain"
)

type fakeDeviceRepo struct {
	devices map[string]domain.Device
}

func newFakeDeviceRepo() *fakeDeviceRepo {
	return &fakeDeviceRepo{devices: map[string]domain.Device{}}
}

func (f *fakeDeviceRepo) Save(d domain.Device) error {
	f.devices[d.ID] = d
	return nil
}

func (f *fakeDeviceRepo) Update(d domain.Device) error {
	f.devices[d.ID] = d
	return nil
}

func (f *fakeDeviceRepo) FindByID(id string) (*domain.Device, error) {
	d, ok := f.devices[id]
	if !ok {
		return nil, nil
	}
	return &d, nil
}

func (f *fakeDeviceRepo) FindBySerial(serial string) (*domain.Device, error) {
	for _, d := range f.devices {
		if d.Serial == serial {
			copy := d
			return &copy, nil
		}
	}
	return nil, nil
}

func (f *fakeDeviceRepo) List() ([]domain.Device, error) {
	out := make([]domain.Device, 0, len(f.devices))
	for _, d := range f.devices {
		out = append(out, d)
	}
	return out, nil
}

type fakeAssetRepo struct {
	assets map[string]domain.Asset
}

func newFakeAssetRepo() *fakeAssetRepo {
	return &fakeAssetRepo{assets: map[string]domain.Asset{}}
}

func (f *fakeAssetRepo) Save(a domain.Asset) error {
	f.assets[a.ID] = a
	return nil
}

func (f *fakeAssetRepo) Update(a domain.Asset) error {
	f.assets[a.ID] = a
	return nil
}

func (f *fakeAssetRepo) FindByID(id string) (*domain.Asset, error) {
	a, ok := f.assets[id]
	if !ok {
		return nil, nil
	}
	return &a, nil
}

func (f *fakeAssetRepo) FindByDeviceID(deviceID string) (*domain.Asset, error) {
	for _, a := range f.assets {
		if a.DeviceID != nil && *a.DeviceID == deviceID {
			copy := a
			return &copy, nil
		}
	}
	return nil, nil
}

func (f *fakeAssetRepo) List() ([]domain.Asset, error) {
	out := make([]domain.Asset, 0, len(f.assets))
	for _, a := range f.assets {
		out = append(out, a)
	}
	return out, nil
}

func TestCreateDeviceGuardaConUUID(t *testing.T) {
	repo := newFakeDeviceRepo()
	uc := NewCreateDevice(repo)

	result, err := uc.Execute(CreateDeviceInput{Serial: "IMEI-001", Name: "Tracker 1"})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if result.Device.ID == "" {
		t.Fatal("esperaba ID generado")
	}
	if result.Device.Protocol != "http" {
		t.Fatalf("protocol = %q, queria http", result.Device.Protocol)
	}
	if result.APIKey == "" {
		t.Fatal("esperaba API key generada")
	}
	if result.Device.APIKeyHash == "" || result.Device.APIKeyHash == result.APIKey {
		t.Fatal("esperaba hash de la API key (no la key en claro)")
	}
	if len(repo.devices) != 1 {
		t.Fatalf("devices = %d, queria 1", len(repo.devices))
	}
}

func TestCreateDeviceSerialDuplicado(t *testing.T) {
	repo := newFakeDeviceRepo()
	uc := NewCreateDevice(repo)

	if _, err := uc.Execute(CreateDeviceInput{Serial: "IMEI-001", Name: "Tracker 1"}); err != nil {
		t.Fatalf("primer alta: %v", err)
	}

	_, err := uc.Execute(CreateDeviceInput{Serial: "IMEI-001", Name: "Tracker 2"})
	if err != domain.ErrSerialDuplicated {
		t.Fatalf("err = %v, queria ErrSerialDuplicated", err)
	}
}

func TestCreateAssetTipoInvalido(t *testing.T) {
	uc := NewCreateAsset(newFakeAssetRepo())

	_, err := uc.Execute(CreateAssetInput{Name: "Camioneta", Type: "camion"})
	if err != domain.ErrInvalidAssetType {
		t.Fatalf("err = %v, queria ErrInvalidAssetType", err)
	}
}

func TestAssignDeviceYaAsignado(t *testing.T) {
	devices := newFakeDeviceRepo()
	assets := newFakeAssetRepo()

	deviceResult, err := NewCreateDevice(devices).Execute(CreateDeviceInput{Serial: "IMEI-100", Name: "Tracker 100"})
	if err != nil {
		t.Fatalf("crear device: %v", err)
	}
	device := deviceResult.Device

	createAsset := NewCreateAsset(assets)
	assetA, err := createAsset.Execute(CreateAssetInput{Name: "Moto A", Type: "moto"})
	if err != nil {
		t.Fatalf("crear asset A: %v", err)
	}
	assetB, err := createAsset.Execute(CreateAssetInput{Name: "Moto B", Type: "moto"})
	if err != nil {
		t.Fatalf("crear asset B: %v", err)
	}

	assign := NewAssignDevice(assets, devices)
	if _, err := assign.Execute(AssignDeviceInput{AssetID: assetA.ID, DeviceID: device.ID}); err != nil {
		t.Fatalf("asignar a A: %v", err)
	}

	_, err = assign.Execute(AssignDeviceInput{AssetID: assetB.ID, DeviceID: device.ID})
	if err != domain.ErrDeviceAlreadyAssigned {
		t.Fatalf("err = %v, queria ErrDeviceAlreadyAssigned", err)
	}

	// Liberar y reasignar si debe funcionar
	if _, err := NewUnassignDevice(assets).Execute(assetA.ID); err != nil {
		t.Fatalf("desasignar A: %v", err)
	}
	if _, err := assign.Execute(AssignDeviceInput{AssetID: assetB.ID, DeviceID: device.ID}); err != nil {
		t.Fatalf("reasignar a B: %v", err)
	}
}

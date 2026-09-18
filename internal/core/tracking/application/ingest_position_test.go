package application

import (
	"testing"
	"time"

	"atlas/internal/core/tracking/domain"
)

type fakeRepo struct {
	saved     []domain.Position
	lastByDev map[string]domain.Position
}

func (f *fakeRepo) Save(p domain.Position) error {
	f.saved = append(f.saved, p)
	if f.lastByDev == nil {
		f.lastByDev = map[string]domain.Position{}
	}
	f.lastByDev[p.DeviceID] = p
	return nil
}

func (f *fakeRepo) LatestByDevice(deviceID string) (*domain.Position, error) {
	p, ok := f.lastByDev[deviceID]
	if !ok {
		return nil, nil
	}
	return &p, nil
}

func (f *fakeRepo) ListByDeviceRange(deviceID string, from, to time.Time, limit int) ([]domain.Position, error) {
	var out []domain.Position
	for _, p := range f.saved {
		if p.DeviceID == deviceID && !p.RecordedAt.Before(from) && p.RecordedAt.Before(to) {
			out = append(out, p)
		}
	}
	return out, nil
}

type fakePublisher struct{ published []domain.Position }

func (f *fakePublisher) PublishPosition(p domain.Position) {
	f.published = append(f.published, p)
}

func TestIngestPositionOK(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePublisher{}
	uc := NewIngestPosition(repo, pub)

	pos, err := uc.Execute(IngestPositionInput{DeviceID: "dev-1", Lat: -12.1, Lng: -77.0})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if pos.ID == "" {
		t.Fatal("esperaba un ID generado")
	}
	if pos.Source != "http" {
		t.Fatalf("source = %q, queria http", pos.Source)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("repo.saved = %d, queria 1", len(repo.saved))
	}
	if len(pub.published) != 1 {
		t.Fatalf("publicados = %d, queria 1", len(pub.published))
	}
}

func TestIngestPositionCoordenadasInvalidas(t *testing.T) {
	uc := NewIngestPosition(&fakeRepo{}, &fakePublisher{})

	_, err := uc.Execute(IngestPositionInput{DeviceID: "dev-1", Lat: 0, Lng: 0})
	if err != domain.ErrInvalidCoordinates {
		t.Fatalf("err = %v, queria ErrInvalidCoordinates", err)
	}
}

func TestIngestPositionDeviceRequerido(t *testing.T) {
	uc := NewIngestPosition(&fakeRepo{}, &fakePublisher{})

	_, err := uc.Execute(IngestPositionInput{Lat: -12.1, Lng: -77.0})
	if err != domain.ErrDeviceRequired {
		t.Fatalf("err = %v, queria ErrDeviceRequired", err)
	}
}

func TestGetHistoryRangoInvalido(t *testing.T) {
	uc := NewGetHistory(&fakeRepo{})

	_, err := uc.Execute(GetHistoryInput{
		DeviceID: "dev-1",
		From:     time.Now(),
		To:       time.Now().Add(-time.Hour),
	})
	if err != ErrInvalidRange {
		t.Fatalf("err = %v, queria ErrInvalidRange", err)
	}
}

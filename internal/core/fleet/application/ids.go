package application

import (
	"strings"

	"github.com/gofrs/uuid/v5"

	"atlas/internal/core/fleet/domain"
)

// parseID valida que un identificador tenga formato UUID (los IDs del modulo
// son UUIDv7). Evita que llegue un string invalido a una columna uuid de
// Postgres, que responderia con error 500.
func parseID(raw string) (string, error) {
	id := strings.TrimSpace(raw)
	if id == "" {
		return "", domain.ErrInvalidID
	}
	if _, err := uuid.FromString(id); err != nil {
		return "", domain.ErrInvalidID
	}
	return id, nil
}

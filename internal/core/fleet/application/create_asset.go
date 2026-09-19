package application

import (
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	"atlas/internal/core/fleet/domain"
)

type CreateAssetInput struct {
	Name  string
	Type  string
	Plate string
}

type CreateAsset struct {
	repo domain.AssetRepository
	now  func() time.Time
}

func NewCreateAsset(repo domain.AssetRepository) *CreateAsset {
	return &CreateAsset{repo: repo, now: time.Now}
}

func (uc *CreateAsset) Execute(input CreateAssetInput) (domain.Asset, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return domain.Asset{}, domain.ErrAssetNameRequired
	}

	assetType, err := domain.ParseAssetType(input.Type)
	if err != nil {
		return domain.Asset{}, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return domain.Asset{}, err
	}

	now := uc.now()
	asset := domain.Asset{
		ID:        id.String(),
		Name:      name,
		Type:      assetType,
		Plate:     strings.ToUpper(strings.TrimSpace(input.Plate)),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.repo.Save(asset); err != nil {
		return domain.Asset{}, err
	}
	return asset, nil
}

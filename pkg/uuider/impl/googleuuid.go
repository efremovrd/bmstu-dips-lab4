package googleuuid

import (
	"bmstu-dips-lab2/pkg/uuider"
	"errors"

	"github.com/google/uuid"
)

type googleuuid struct{}

func NewGoogleUUID() uuider.UUIDer {
	return &googleuuid{}
}

func (g *googleuuid) Generate() (*string, error) {
	generated := uuid.New()
	converted := generated.String()
	if converted == "" {
		return nil, errors.New("invalid generated uuid")
	}

	return &converted, nil
}

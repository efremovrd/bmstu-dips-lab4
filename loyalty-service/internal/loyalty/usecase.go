package loyalty

import (
	"bmstu-dips-lab2/loyalty-service/models"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, model *models.Loyalty) (*models.Loyalty, error)
	UpdateResCountByOne(ctx context.Context, model *models.Loyalty) (*models.Loyalty, error)
	GetByUsername(ctx context.Context, username string) (*models.Loyalty, error)
}

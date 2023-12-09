package loyalty

import (
	"bmstu-dips-lab2/loyalty-service/models"
	"context"
)

type Repo interface {
	Create(ctx context.Context, model *models.Loyalty) (*models.Loyalty, error)
	GetByUsername(ctx context.Context, username string) (*models.Loyalty, error)
	Update(ctx context.Context, model *models.Loyalty, toUpdate *models.Loyalty) (*models.Loyalty, error)
}

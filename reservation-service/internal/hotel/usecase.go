package hotel

import (
	"bmstu-dips-lab2/pkg/types"
	"bmstu-dips-lab2/reservation-service/models"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, model *models.Hotel) (*models.Hotel, error)
	GetAllPaged(ctx context.Context, sets types.GetSets) ([]*models.Hotel, error)
	GetById(ctx context.Context, id int) (*models.Hotel, error)
	GetByUid(ctx context.Context, uid string) (*models.Hotel, error)
}

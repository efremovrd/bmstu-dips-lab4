package reservation

import (
	"bmstu-dips-lab2/reservation-service/models"
	"context"
)

type UseCase interface {
	Create(ctx context.Context, model *models.Reservation) (*models.Reservation, error)
	GetByUsername(ctx context.Context, username string) ([]*models.Reservation, error)
	GetByReservationUid(ctx context.Context, reservation_uid, username string) (*models.Reservation, error)
	Delete(ctx context.Context, reservation_uid string) error
}
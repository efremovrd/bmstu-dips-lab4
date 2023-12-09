package usecase

import (
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/pkg/uuider"
	"bmstu-dips-lab2/reservation-service/internal/reservation"
	"bmstu-dips-lab2/reservation-service/models"
	"context"
)

type ReservationUseCase struct {
	reservationRepo reservation.Repo
	uuider          uuider.UUIDer
}

func NewReservationUseCase(reservationRepo reservation.Repo, uuider uuider.UUIDer) reservation.UseCase {
	return &ReservationUseCase{
		reservationRepo: reservationRepo,
		uuider:          uuider,
	}
}

func (r *ReservationUseCase) Create(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error) {
	if reservation.Reservation_uid == "" {
		newUUID, err := r.uuider.Generate()
		if err != nil {
			return nil, err
		}

		reservation.Reservation_uid = *newUUID
	}

	return r.reservationRepo.Create(ctx, reservation)
}

func (r *ReservationUseCase) GetByUsername(ctx context.Context, username string) ([]*models.Reservation, error) {
	return r.reservationRepo.GetByUsername(ctx, username)
}

func (r *ReservationUseCase) GetByReservationUid(ctx context.Context, reservation_uid, username string) (*models.Reservation, error) {
	found, err := r.reservationRepo.GetByReservationUid(ctx, reservation_uid)
	if err != nil {
		return nil, err
	}

	if found.Username != username {
		return nil, errs.ErrForbidden
	}

	return found, nil
}

func (r *ReservationUseCase) Delete(ctx context.Context, reservation_uid string) error {
	return r.reservationRepo.Delete(ctx, reservation_uid)
}

package usecase

import (
	"bmstu-dips-lab2/pkg/types"
	"bmstu-dips-lab2/pkg/uuider"
	"bmstu-dips-lab2/reservation-service/internal/hotel"
	"bmstu-dips-lab2/reservation-service/models"
	"context"
)

type HotelUseCase struct {
	hotelRepo hotel.Repo
	uuider    uuider.UUIDer
}

func NewHotelUseCase(hotelRepo hotel.Repo, uuider uuider.UUIDer) hotel.UseCase {
	return &HotelUseCase{
		hotelRepo: hotelRepo,
		uuider:    uuider,
	}
}

func (h *HotelUseCase) Create(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	if hotel.Hotel_uid == "" {
		newUUID, err := h.uuider.Generate()
		if err != nil {
			return nil, err
		}

		hotel.Hotel_uid = *newUUID
	}

	return h.hotelRepo.Create(ctx, hotel)
}

func (h *HotelUseCase) GetAllPaged(ctx context.Context, sets types.GetSets) ([]*models.Hotel, error) {
	return h.hotelRepo.GetAllPaged(ctx, sets)
}

func (h *HotelUseCase) GetById(ctx context.Context, id int) (*models.Hotel, error) {
	return h.hotelRepo.GetById(ctx, id)
}

func (h *HotelUseCase) GetByUid(ctx context.Context, uid string) (*models.Hotel, error) {
	return h.hotelRepo.GetByUid(ctx, uid)
}

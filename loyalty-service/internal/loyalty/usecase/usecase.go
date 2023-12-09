package usecase

import (
	"bmstu-dips-lab2/loyalty-service/internal/loyalty"
	"bmstu-dips-lab2/loyalty-service/models"
	"bmstu-dips-lab2/pkg/errs"
	"context"
)

type LoyaltyUseCase struct {
	loyaltyRepo loyalty.Repo
}

func NewLoyaltyUseCase(loyaltyRepo loyalty.Repo) loyalty.UseCase {
	return &LoyaltyUseCase{
		loyaltyRepo: loyaltyRepo,
	}
}

func (l *LoyaltyUseCase) Create(ctx context.Context, loyalty *models.Loyalty) (*models.Loyalty, error) {
	defauldiscount, err := GetDiscount("BRONZE")
	if err != nil {
		return nil, err
	}
	loyalty.Discount = defauldiscount
	return l.loyaltyRepo.Create(ctx, loyalty)
}

func (l *LoyaltyUseCase) UpdateResCountByOne(ctx context.Context, loyalty *models.Loyalty) (*models.Loyalty, error) {
	if loyalty.Username == "" || loyalty.Reservation_count == 0 {
		return nil, errs.ErrInvalidContent
	}

	foundloyalty, err := l.loyaltyRepo.GetByUsername(ctx, loyalty.Username)
	if err != nil {
		return nil, err
	}

	toUpdate := new(models.Loyalty)

	if loyalty.Reservation_count > 0 {
		foundloyalty.Reservation_count = foundloyalty.Reservation_count + 1
	} else {
		if foundloyalty.Reservation_count == 0 {
			return foundloyalty, nil
		}
		foundloyalty.Reservation_count = foundloyalty.Reservation_count - 1
	}
	toUpdate.Reservation_count = 1

	newstatus, err := GetStatus(foundloyalty.Reservation_count)
	if err != nil {
		return nil, err
	}

	if foundloyalty.Status != newstatus {
		foundloyalty.Status = newstatus
		toUpdate.Status = "y"

		newdiscount, err := GetDiscount(newstatus)
		if err != nil {
			return nil, err
		}
		foundloyalty.Discount = newdiscount
		toUpdate.Discount = 1

	}

	return l.loyaltyRepo.Update(ctx, foundloyalty, toUpdate)
}

func (l *LoyaltyUseCase) GetByUsername(ctx context.Context, username string) (*models.Loyalty, error) {
	return l.loyaltyRepo.GetByUsername(ctx, username)
}

func GetStatus(reservation_count int) (string, error) {
	if reservation_count < 0 {
		return "", errs.ErrInvalidContent
	}

	if reservation_count < 10 {
		return "BRONZE", nil
	}

	if reservation_count < 20 {
		return "SILVER", nil
	}

	return "GOLD", nil
}

func GetDiscount(status string) (int, error) {
	if status == "BRONZE" {
		return 5, nil
	}

	if status == "SILVER" {
		return 7, nil
	}

	if status == "GOLD" {
		return 10, nil
	}

	return 0, errs.ErrInvalidContent
}

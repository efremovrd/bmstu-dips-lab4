package repo

import (
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/pkg/postgres"
	"bmstu-dips-lab2/pkg/types"
	"bmstu-dips-lab2/reservation-service/internal/hotel"
	"bmstu-dips-lab2/reservation-service/models"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

type HotelDB struct {
	id, stars, price                        int
	name, country, address, city, hotel_uid string
}

type HotelRepo struct {
	*postgres.Postgres
}

func NewHotelRepo(db *postgres.Postgres) hotel.Repo {
	return &HotelRepo{db}
}

func (h *HotelRepo) Create(ctx context.Context, modelBL *models.Hotel) (*models.Hotel, error) {
	modelDB, err := HotelBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := h.Builder.
		Insert("hotels").
		Columns("price, stars, hotel_uid, city, address, country, name").
		Values(modelDB.price, modelDB.stars, modelDB.hotel_uid, modelDB.city, modelDB.address, modelDB.country, modelDB.name).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = h.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id)
	if err != nil {
		return nil, err
	}

	return HotelDBToBL(modelDB)
}

func (h *HotelRepo) GetAllPaged(ctx context.Context, sets types.GetSets) ([]*models.Hotel, error) {
	sql, args, err := h.Builder.
		Select("id, price, stars, hotel_uid, city, address, country, name").
		From("hotels").
		Limit(sets.Limit).
		Offset(sets.Offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := h.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Hotel, 0)

	for rows.Next() {
		modelDB := HotelDB{}

		err = rows.Scan(&modelDB.id, &modelDB.price, &modelDB.stars, &modelDB.hotel_uid, &modelDB.city, &modelDB.address, &modelDB.country, &modelDB.name)
		if err != nil {
			return nil, err
		}

		questionBL, err := HotelDBToBL(&modelDB)
		if err != nil {
			return nil, err
		}

		res = append(res, questionBL)
	}

	return res, nil
}

func (h *HotelRepo) GetById(ctx context.Context, id int) (*models.Hotel, error) {
	sql, args, err := h.Builder.
		Select("price, stars, hotel_uid, city, address, country, name").
		From("hotels").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := HotelDB{id: id}
	err = h.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.price, &modelDB.stars, &modelDB.hotel_uid, &modelDB.city, &modelDB.address, &modelDB.country, &modelDB.name)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return HotelDBToBL(&modelDB)
}

func (h *HotelRepo) GetByUid(ctx context.Context, uid string) (*models.Hotel, error) {
	sql, args, err := h.Builder.
		Select("price, stars, id, city, address, country, name").
		From("hotels").
		Where(squirrel.Eq{"hotel_uid": uid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := HotelDB{hotel_uid: uid}
	err = h.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.price, &modelDB.stars, &modelDB.id, &modelDB.city, &modelDB.address, &modelDB.country, &modelDB.name)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return HotelDBToBL(&modelDB)
}

func HotelDBToBL(modelDB *HotelDB) (*models.Hotel, error) {
	return &models.Hotel{
		Id:        modelDB.id,
		Stars:     modelDB.stars,
		Price:     modelDB.price,
		Hotel_uid: modelDB.hotel_uid,
		Address:   modelDB.address,
		Country:   modelDB.country,
		Name:      modelDB.name,
		City:      modelDB.city,
	}, nil
}

func HotelBLToDB(modelBL *models.Hotel) (*HotelDB, error) {
	return &HotelDB{
		id:        modelBL.Id,
		stars:     modelBL.Stars,
		price:     modelBL.Price,
		hotel_uid: modelBL.Hotel_uid,
		address:   modelBL.Address,
		country:   modelBL.Country,
		name:      modelBL.Name,
		city:      modelBL.City,
	}, nil
}

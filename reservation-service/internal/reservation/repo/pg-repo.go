package repo

import (
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/pkg/postgres"
	"bmstu-dips-lab2/reservation-service/internal/reservation"
	"bmstu-dips-lab2/reservation-service/models"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

type ReservationDB struct {
	id, hotel_id                                                         int
	username, status, reservation_uid, payment_uid, start_date, end_data string
}

type ReservationRepo struct {
	*postgres.Postgres
}

func NewReservationRepo(db *postgres.Postgres) reservation.Repo {
	return &ReservationRepo{db}
}

func (r *ReservationRepo) Create(ctx context.Context, modelBL *models.Reservation) (*models.Reservation, error) {
	modelDB, err := ReservationBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := r.Builder.
		Insert("reservation").
		Columns("hotel_id, username, status, reservation_uid, payment_uid, start_date, end_data").
		Values(modelDB.hotel_id, modelDB.username, modelDB.status, modelDB.reservation_uid, modelDB.payment_uid, modelDB.start_date, modelDB.end_data).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id)
	if err != nil {
		return nil, err
	}

	return ReservationDBToBL(modelDB)
}

func (r *ReservationRepo) GetByUsername(ctx context.Context, username string) ([]*models.Reservation, error) {
	sql, args, err := r.Builder.
		Select("id, hotel_id, status, reservation_uid, payment_uid, start_date::text, end_data::text").
		From("reservation").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Reservation, 0)

	for rows.Next() {
		modelDB := ReservationDB{username: username}

		err = rows.Scan(&modelDB.id, &modelDB.hotel_id, &modelDB.status, &modelDB.reservation_uid, &modelDB.payment_uid, &modelDB.start_date, &modelDB.end_data)
		if err != nil {
			return nil, err
		}

		reservationBL, err := ReservationDBToBL(&modelDB)
		if err != nil {
			return nil, err
		}

		res = append(res, reservationBL)
	}

	return res, nil
}

func (r *ReservationRepo) GetByReservationUid(ctx context.Context, reservation_uid string) (*models.Reservation, error) {
	sql, args, err := r.Builder.
		Select("id, hotel_id, status, username, payment_uid, start_date::text, end_data::text").
		From("reservation").
		Where(squirrel.Eq{"reservation_uid": reservation_uid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := ReservationDB{reservation_uid: reservation_uid}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id, &modelDB.hotel_id, &modelDB.status, &modelDB.username, &modelDB.payment_uid, &modelDB.start_date, &modelDB.end_data)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return ReservationDBToBL(&modelDB)
}

func (r *ReservationRepo) Delete(ctx context.Context, reservation_uid string) error {
	sql, args, err := r.Builder.
		Update("reservation").
		Set("status", "CANCELED").
		Where(squirrel.Eq{"reservation_uid": reservation_uid}).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func ReservationBLToDB(modelBL *models.Reservation) (*ReservationDB, error) {
	return &ReservationDB{
		id:              modelBL.Id,
		status:          modelBL.Status,
		username:        modelBL.Username,
		hotel_id:        modelBL.Hotel_id,
		reservation_uid: modelBL.Reservation_uid,
		payment_uid:     modelBL.Payment_uid,
		start_date:      modelBL.Start_date,
		end_data:        modelBL.End_data,
	}, nil
}

func ReservationDBToBL(modelDB *ReservationDB) (*models.Reservation, error) {
	return &models.Reservation{
		Id:              modelDB.id,
		Status:          modelDB.status,
		Username:        modelDB.username,
		Hotel_id:        modelDB.hotel_id,
		Reservation_uid: modelDB.reservation_uid,
		Payment_uid:     modelDB.payment_uid,
		Start_date:      modelDB.start_date,
		End_data:        modelDB.end_data,
	}, nil
}

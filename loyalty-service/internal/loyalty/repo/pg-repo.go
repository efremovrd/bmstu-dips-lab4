package repo

import (
	"bmstu-dips-lab2/loyalty-service/internal/loyalty"
	"bmstu-dips-lab2/loyalty-service/models"
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/pkg/postgres"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

type LoyaltyDB struct {
	id, discount, reservation_count int
	status, username                string
}

type LoyaltyRepo struct {
	*postgres.Postgres
}

func NewLoyaltyRepo(db *postgres.Postgres) loyalty.Repo {
	return &LoyaltyRepo{db}
}

func (l *LoyaltyRepo) Create(ctx context.Context, modelBL *models.Loyalty) (*models.Loyalty, error) {
	modelDB, err := LoyaltyBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := l.Builder.
		Insert("loyalty").
		Columns("username, discount").
		Values(modelDB.username, modelDB.discount).
		Suffix("RETURNING id, status, reservation_count").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = l.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id, &modelDB.status, &modelDB.reservation_count)
	if err != nil {
		return nil, err
	}

	return LoyaltyDBToBL(modelDB)
}

func (l *LoyaltyRepo) Update(ctx context.Context, modelBL *models.Loyalty, toUpdate *models.Loyalty) (*models.Loyalty, error) {
	modelDB, err := LoyaltyBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	builder := l.Builder.
		Update("loyalty")

	if toUpdate.Status != "" {
		builder = builder.
			Set("status", modelDB.status)
	}

	if toUpdate.Reservation_count != 0 {
		builder = builder.
			Set("reservation_count", modelDB.reservation_count)
	}

	if toUpdate.Discount != 0 {
		builder = builder.
			Set("discount", modelDB.discount)
	}

	sql, args, err := builder.
		Where(squirrel.Eq{"username": modelDB.username}).
		Suffix("RETURNING status, reservation_count, discount, id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = l.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.status, &modelDB.reservation_count, &modelDB.discount, &modelDB.id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errs.ErrNoContent
		}
		return nil, err
	}

	return LoyaltyDBToBL(modelDB)
}

func (l *LoyaltyRepo) GetByUsername(ctx context.Context, username string) (*models.Loyalty, error) {
	sql, args, err := l.Builder.
		Select("id, status, discount, reservation_count").
		From("loyalty").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := LoyaltyDB{username: username}
	err = l.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id, &modelDB.status, &modelDB.discount, &modelDB.reservation_count)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return LoyaltyDBToBL(&modelDB)
}

func LoyaltyDBToBL(modelDB *LoyaltyDB) (*models.Loyalty, error) {
	return &models.Loyalty{
		Id:                modelDB.id,
		Status:            modelDB.status,
		Reservation_count: modelDB.reservation_count,
		Username:          modelDB.username,
		Discount:          modelDB.discount,
	}, nil
}

func LoyaltyBLToDB(modelBL *models.Loyalty) (*LoyaltyDB, error) {
	return &LoyaltyDB{
		id:                modelBL.Id,
		status:            modelBL.Status,
		reservation_count: modelBL.Reservation_count,
		username:          modelBL.Username,
		discount:          modelBL.Discount,
	}, nil
}

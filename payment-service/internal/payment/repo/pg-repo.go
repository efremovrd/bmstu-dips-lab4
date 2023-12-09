package repo

import (
	"bmstu-dips-lab2/payment-service/internal/payment"
	"bmstu-dips-lab2/payment-service/models"
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/pkg/postgres"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
)

type PaymentDB struct {
	id, price           int
	status, payment_uid string
}

type PaymentRepo struct {
	*postgres.Postgres
}

func NewPaymentRepo(db *postgres.Postgres) payment.Repo {
	return &PaymentRepo{db}
}

func (p *PaymentRepo) Create(ctx context.Context, modelBL *models.Payment) (*models.Payment, error) {
	modelDB, err := PaymentBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	sql, args, err := p.Builder.
		Insert("payment").
		Columns("price, status, payment_uid").
		Values(modelDB.price, modelDB.status, modelDB.payment_uid).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id)
	if err != nil {
		return nil, err
	}

	return PaymentDBToBL(modelDB)
}

func (p *PaymentRepo) Update(ctx context.Context, modelBL *models.Payment, toUpdate *models.Payment) (*models.Payment, error) {
	modelDB, err := PaymentBLToDB(modelBL)
	if err != nil {
		return nil, errs.ErrInvalidContent
	}

	builder := p.Builder.
		Update("payment")

	if toUpdate.Status != "" {
		builder = builder.
			Set("status", modelDB.status)
	}

	if toUpdate.Price != 0 {
		builder = builder.
			Set("price", modelDB.price)
	}

	sql, args, err := builder.
		Where(squirrel.Eq{"payment_uid": modelDB.payment_uid}).
		Suffix("RETURNING status, price, id").
		ToSql()
	if err != nil {
		return nil, err
	}

	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.status, &modelDB.price, &modelDB.id)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return PaymentDBToBL(modelDB)
}

func (p *PaymentRepo) GetByPaymentUid(ctx context.Context, payment_uid string) (*models.Payment, error) {
	sql, args, err := p.Builder.
		Select("id, price, status").
		From("payment").
		Where(squirrel.Eq{"payment_uid": payment_uid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	modelDB := PaymentDB{payment_uid: payment_uid}
	err = p.Pool.QueryRow(ctx, sql, args...).Scan(&modelDB.id, &modelDB.price, &modelDB.status)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return PaymentDBToBL(&modelDB)
}

func PaymentDBToBL(modelDB *PaymentDB) (*models.Payment, error) {
	return &models.Payment{
		Id:          modelDB.id,
		Status:      modelDB.status,
		Price:       modelDB.price,
		Payment_uid: modelDB.payment_uid,
	}, nil
}

func PaymentBLToDB(modelBL *models.Payment) (*PaymentDB, error) {
	return &PaymentDB{
		id:          modelBL.Id,
		status:      modelBL.Status,
		price:       modelBL.Price,
		payment_uid: modelBL.Payment_uid,
	}, nil
}

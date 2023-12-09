package repo_test

import (
	"bmstu-dips-lab2/payment-service/internal/payment/repo"
	"bmstu-dips-lab2/payment-service/models"
	"bmstu-dips-lab2/pkg/postgres"
	"context"
	"errors"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	_builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func TestPaymentRepo_GetByPaymentUid(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewPaymentRepo(&db)

	type mockBehavior func(ctx context.Context, uid string)

	testTable := []struct {
		nameTest            string
		ctx                 context.Context
		uid                 string
		reservation         models.Payment
		mockBehavior        mockBehavior
		expectedReservation models.Payment
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			uid:      "someUid",
			mockBehavior: func(ctx context.Context, uid string) {
				pgxRows := pgxpoolmock.NewRows([]string{"id", "price", "status"}).AddRow(1, 2, "a").ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id, price, status FROM payment WHERE payment_uid = $1", uid).Return(pgxRows)
			},
			expectedReservation: models.Payment{
				Id:          1,
				Price:       2,
				Status:      "a",
				Payment_uid: "someUid",
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			uid:      "someUid",
			mockBehavior: func(ctx context.Context, uid string) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id, price, status FROM payment WHERE payment_uid = $1", uid).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.uid)

			got, err := r.GetByPaymentUid(testCase.ctx, testCase.uid)

			switch testCase.nameTest {
			case "ok":
				assert.Equal(t, nil, err)
				assert.Equal(t, testCase.expectedReservation, *got)
			case "no_rows":
				assert.NotEqual(t, nil, err)
			default:
				assert.Error(t, errors.New("No case"), "No case")
			}
		})
	}
}

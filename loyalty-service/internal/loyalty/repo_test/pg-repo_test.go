package repo_test

import (
	"bmstu-dips-lab2/loyalty-service/internal/loyalty/repo"
	"bmstu-dips-lab2/loyalty-service/models"
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

func TestLoyaltyRepo_GetByPaymentUid(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewLoyaltyRepo(&db)

	type mockBehavior func(ctx context.Context, username string)

	testTable := []struct {
		nameTest            string
		ctx                 context.Context
		username            string
		reservation         models.Loyalty
		mockBehavior        mockBehavior
		expectedReservation models.Loyalty
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			username: "someUsername",
			mockBehavior: func(ctx context.Context, uid string) {
				pgxRows := pgxpoolmock.NewRows([]string{"id", "status", "disount", "reservation_count"}).AddRow(1, "a", 3, 2).ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id, status, discount, reservation_count FROM loyalty WHERE username = $1", uid).Return(pgxRows)
			},
			expectedReservation: models.Loyalty{
				Id:                1,
				Reservation_count: 2,
				Discount:          3,
				Status:            "a",
				Username:          "someUsername",
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			username: "someUsername",
			mockBehavior: func(ctx context.Context, uid string) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id, status, discount, reservation_count FROM loyalty WHERE username = $1", uid).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.username)

			got, err := r.GetByUsername(testCase.ctx, testCase.username)

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

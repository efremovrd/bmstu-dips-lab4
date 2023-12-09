package repo_test

import (
	"bmstu-dips-lab2/pkg/postgres"
	"bmstu-dips-lab2/reservation-service/internal/reservation/repo"
	"bmstu-dips-lab2/reservation-service/models"
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

func TestReservationRepo_GetByReservationUid(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

	db := postgres.Postgres{
		Builder: _builder,
		Pool:    mockPool,
	}

	r := repo.NewReservationRepo(&db)

	type mockBehavior func(ctx context.Context, uid string)

	testTable := []struct {
		nameTest            string
		ctx                 context.Context
		uid                 string
		reservation         models.Reservation
		mockBehavior        mockBehavior
		expectedReservation models.Reservation
	}{
		{
			nameTest: "ok",
			ctx:      context.Background(),
			uid:      "someUid",
			mockBehavior: func(ctx context.Context, uid string) {
				pgxRows := pgxpoolmock.NewRows([]string{"id", "hotel_id", "status", "username", "payment_uid", "start_date", "end_data"}).AddRow(1, 2, "a", "b", "c", "d", "e").ToPgxRows()
				pgxRows.Next()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id, hotel_id, status, username, payment_uid, start_date::text, end_data::text FROM reservation WHERE reservation_uid = $1", uid).Return(pgxRows)
			},
			expectedReservation: models.Reservation{
				Id:              1,
				Hotel_id:        2,
				Status:          "a",
				Username:        "b",
				Payment_uid:     "c",
				Start_date:      "d",
				End_data:        "e",
				Reservation_uid: "someUid",
			},
		},
		{
			nameTest: "no_rows",
			ctx:      context.Background(),
			uid:      "someUid",
			mockBehavior: func(ctx context.Context, uid string) {
				pgxRows := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()
				mockPool.EXPECT().QueryRow(ctx, "SELECT id, hotel_id, status, username, payment_uid, start_date::text, end_data::text FROM reservation WHERE reservation_uid = $1", uid).Return(pgxRows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.nameTest, func(t *testing.T) {
			testCase.mockBehavior(testCase.ctx, testCase.uid)

			got, err := r.GetByReservationUid(testCase.ctx, testCase.uid)

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

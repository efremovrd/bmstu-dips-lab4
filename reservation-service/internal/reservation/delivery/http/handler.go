package http

import (
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/reservation-service/internal/hotel"
	"bmstu-dips-lab2/reservation-service/internal/reservation"
	"bmstu-dips-lab2/reservation-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReservationCreatRequest struct {
	Payment_uid string `json:"paymentUid" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Hotel_id    int    `json:"hotel_id" binding:"required"`
	Start_date  string `json:"startDate" binding:"required"`
	End_data    string `json:"endDate" binding:"required"`
}

type HotelResponse struct {
	Stars       int    `json:"stars"`
	FullAddress string `json:"fullAddress"`
	Name        string `json:"name"`
	Hotel_uid   string `json:"hotelUid"`
}

type ReservationResponse struct {
	Id              int           `json:"id"`
	Hotel_id        int           `json:"hotel_id"`
	Payment_uid     string        `json:"paymentUid"`
	Reservation_uid string        `json:"reservationUid"`
	Username        string        `json:"username"`
	Status          string        `json:"status"`
	Start_date      string        `json:"startDate"`
	End_data        string        `json:"endDate"`
	Hotel           HotelResponse `json:"hotel"`
}

type ReservationHandlers struct {
	reservationUC reservation.UseCase
	hotelUC       hotel.UseCase
}

func NewReservationHandlers(reservationUC reservation.UseCase, hotelUC hotel.UseCase) reservation.Handlers {
	return &ReservationHandlers{
		reservationUC: reservationUC,
		hotelUC:       hotelUC,
	}
}

func (r *ReservationHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(ReservationCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := ReservationCreatRequestToBL(request)

		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		modelBL.Username = username

		createdReservation, err := r.reservationUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Header("Location", "/api/v1/reservations/"+createdReservation.Reservation_uid)

		c.Status(http.StatusCreated)
	}
}

func (r *ReservationHandlers) GetByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundreservations, err := r.reservationUC.GetByUsername(c, username)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		foundhotels := make([]*models.Hotel, len(foundreservations))
		for i, reserv := range foundreservations {
			foundhotel, _ := r.hotelUC.GetById(c, reserv.Hotel_id)
			foundhotels[i] = foundhotel
		}

		c.JSON(http.StatusOK, ReservationsBLToResponse(foundreservations, foundhotels))
	}
}

func (r *ReservationHandlers) GetByReservationUid() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		reservation_uid := c.Param("reservationUid")
		if reservation_uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundreservation, err := r.reservationUC.GetByReservationUid(c, reservation_uid, username)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		foundhotel, err := r.hotelUC.GetById(c, foundreservation.Hotel_id)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, ReservationBLToResponse(foundreservation, foundhotel))
	}
}

func (r *ReservationHandlers) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		reservation_uid := c.Param("reservationUid")
		if reservation_uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err := r.reservationUC.Delete(c, reservation_uid)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func ReservationCreatRequestToBL(dto *ReservationCreatRequest) *models.Reservation {
	return &models.Reservation{
		Status:      dto.Status,
		Start_date:  dto.Start_date,
		End_data:    dto.End_data,
		Hotel_id:    dto.Hotel_id,
		Payment_uid: dto.Payment_uid,
	}
}

func ReservationBLToResponse(model *models.Reservation, hotel *models.Hotel) *ReservationResponse {
	return &ReservationResponse{
		Id:              model.Id,
		Hotel_id:        model.Hotel_id,
		Payment_uid:     model.Payment_uid,
		Reservation_uid: model.Reservation_uid,
		Username:        model.Username,
		Status:          model.Status,
		Start_date:      model.Start_date,
		End_data:        model.End_data,
		Hotel: HotelResponse{
			Hotel_uid:   hotel.Hotel_uid,
			Stars:       hotel.Stars,
			Name:        hotel.Name,
			FullAddress: hotel.Country + ", " + hotel.City + ", " + hotel.Address,
		},
	}
}

func ReservationsBLToResponse(reservations []*models.Reservation, hotels []*models.Hotel) []*ReservationResponse {
	if reservations == nil {
		return nil
	}

	res := make([]*ReservationResponse, len(reservations))

	for i, p := range reservations {

		res[i] = ReservationBLToResponse(p, hotels[i])
	}

	return res
}

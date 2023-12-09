package http

import (
	"bmstu-dips-lab2/gateway-service/internal/gateway"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	loylatyService     string = "http://loyalty-service:8050"
	paymentService     string = "http://payment-service:8060"
	reservationService string = "http://reservation-service:8070"
)

type LoyaltyResponse struct {
	Status            string `json:"status"`
	Reservation_count int    `json:"reservationCount"`
	Discount          int    `json:"discount"`
}

type HotelResponse struct {
	Id        int    `json:"id"`
	Price     int    `json:"price"`
	Stars     int    `json:"stars"`
	Name      string `json:"name"`
	Hotel_uid string `json:"hotelUid"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Country   string `json:"country"`
}

type HotelsResponse struct {
	Page          int              `json:"page"`
	PageSize      int              `json:"pageSize"`
	TotalElements int              `json:"totalElements"`
	Items         []*HotelResponse `json:"items"`
}

type PostRequest struct {
	Hotel_uid  string `json:"hotelUid" binding:"required"`
	Start_date string `json:"startDate" binding:"required"`
	End_date   string `json:"endDate" binding:"required"`
}

type PaymentResponse struct {
	Price  int    `json:"price"`
	Status string `json:"status"`
}

type CreateReservationResponse struct {
	Reservation_uid string           `json:"reservationUid"`
	Hotel_uid       string           `json:"hotelUid"`
	Start_date      string           `json:"startDate"`
	End_date        string           `json:"endDate"`
	Discount        int              `json:"discount"`
	Status          string           `json:"status"`
	Payment         *PaymentResponse `json:"payment"`
}

type LoyaltyCreatRequest struct {
	Username string `json:"username" binding:"required"`
}

type HotelFAResponse struct {
	Stars       int    `json:"stars"`
	FullAddress string `json:"fullAddress"`
	Name        string `json:"name"`
	Hotel_uid   string `json:"hotelUid"`
}

type GetReservationByUid struct {
	Id              int              `json:"id"`
	Hotel_id        int              `json:"hotel_id"`
	Payment_uid     string           `json:"paymentUid"`
	Reservation_uid string           `json:"reservationUid"`
	Username        string           `json:"username"`
	Status          string           `json:"status"`
	Start_date      string           `json:"startDate"`
	End_data        string           `json:"endDate"`
	Hotel           *HotelFAResponse `json:"hotel"`
}

type GetOneReservation struct {
	Reservation_uid string           `json:"reservationUid"`
	Hotel           HotelFAResponse  `json:"hotel"`
	Status          string           `json:"status"`
	Start_date      string           `json:"startDate"`
	End_data        string           `json:"endDate"`
	Payment         *PaymentResponse `json:"payment"`
}

type UserLoyalty struct {
	Status   string `json:"status"`
	Discount int    `json:"discount"`
}

type UserInfo struct {
	Reservations []*GetOneReservation `json:"reservations"`
	Loyalty      *UserLoyalty         `json:"loyalty"`
}

type GatewayHandlers struct{}

func NewGatewayHandlers() gateway.Handlers {
	return &GatewayHandlers{}
}

func (g *GatewayHandlers) GetLoyalty() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		requestURL := loylatyService + c.FullPath()
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		req.Header.Set("X-User-Name", username)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if res.StatusCode != http.StatusOK {
			c.AbortWithStatus(res.StatusCode)
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var data LoyaltyResponse
		if json.Unmarshal(resBody, &data) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(res.StatusCode, data)
	}
}

func (g *GatewayHandlers) GetHotels() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURL := reservationService + c.FullPath() + "?page=" + c.Query("page") + "&size=" + c.Query("size")
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if res.StatusCode != http.StatusOK {
			c.AbortWithStatus(res.StatusCode)
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var data HotelsResponse
		if json.Unmarshal(resBody, &data) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(res.StatusCode, &data)
	}
}

func (g *GatewayHandlers) CreateReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")

		ownrequest := new(PostRequest)

		err := c.ShouldBindJSON(ownrequest)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		checkhotelurl := reservationService + "/api/v1/hotels/" + ownrequest.Hotel_uid
		checkhotelreq, err := http.NewRequest(http.MethodGet, checkhotelurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		checkhotelres, err := http.DefaultClient.Do(checkhotelreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if checkhotelres.StatusCode != http.StatusOK {
			c.AbortWithStatus(checkhotelres.StatusCode)
			return
		}

		checkhotelresresBody, err := io.ReadAll(checkhotelres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var checkhotel HotelResponse
		if json.Unmarshal(checkhotelresresBody, &checkhotel) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		layout := "2006-01-02" // "2006-01-02T15:04:05.000-03:00"
		startdate, err := time.Parse(layout, ownrequest.Start_date)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		enddate, err := time.Parse(layout, ownrequest.End_date)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		days := int(enddate.Sub(startdate).Hours() / 24)
		price := days * checkhotel.Price

		checkloyaltyurl := loylatyService + "/api/v1/loyalty"
		checkloyaltyreq, err := http.NewRequest(http.MethodGet, checkloyaltyurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		checkloyaltyreq.Header.Set("X-User-Name", username)

		checkloyaltyres, err := http.DefaultClient.Do(checkloyaltyreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if checkloyaltyres.StatusCode == http.StatusNotFound {
			jsonBody := []byte("{\"username\": \"" + username + "\"}")
			bodyReader := bytes.NewReader(jsonBody)

			createloyaltyreq, err := http.NewRequest(http.MethodPost, checkloyaltyurl, bodyReader)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			createloyaltyres, err := http.DefaultClient.Do(createloyaltyreq)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if createloyaltyres.StatusCode != http.StatusCreated {
				c.AbortWithStatus(createloyaltyres.StatusCode)
				return
			}
		} else if checkloyaltyres.StatusCode != http.StatusOK {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		checkloyaltyres, _ = http.DefaultClient.Do(checkloyaltyreq)

		checkloyaltyresBody, err := io.ReadAll(checkloyaltyres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var checkloyaltydata LoyaltyResponse
		if json.Unmarshal(checkloyaltyresBody, &checkloyaltydata) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		price = int(float32(price) * (1 - float32(checkloyaltydata.Discount)/100))

		createpaymentjsonBody := []byte(fmt.Sprintf("{\"price\": %d, \"status\": \"PAID\"}", price))
		createpaymentbodyReader := bytes.NewReader(createpaymentjsonBody)

		createpaymenturl := paymentService + "/api/v1/payments"
		createpaymentreq, err := http.NewRequest(http.MethodPost, createpaymenturl, createpaymentbodyReader)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		createpaymentres, err := http.DefaultClient.Do(createpaymentreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if createpaymentres.StatusCode != http.StatusCreated {
			c.AbortWithStatus(createpaymentres.StatusCode)
			return
		}

		temp := strings.Split(createpaymentres.Header.Get("Location"), "/")
		payment_uid := temp[len(temp)-1]

		// increment reservation count
		incrementcounterjsonBody := []byte("{\"reservationCount\": 1}")
		incrementcounterbodyReader := bytes.NewReader(incrementcounterjsonBody)

		incrementcounterurl := loylatyService + "/api/v1/loyalty"
		incrementcounterreq, err := http.NewRequest(http.MethodPatch, incrementcounterurl, incrementcounterbodyReader)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		incrementcounterreq.Header.Set("X-User-Name", username)

		incrementcounterres, err := http.DefaultClient.Do(incrementcounterreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if incrementcounterres.StatusCode != http.StatusOK {
			c.AbortWithStatus(incrementcounterres.StatusCode)
			return
		}

		incrementcounterresBody, err := io.ReadAll(incrementcounterres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var incrementcounterdata LoyaltyResponse
		if json.Unmarshal(incrementcounterresBody, &incrementcounterdata) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// create reservation
		createreservationjsonBody := []byte(fmt.Sprintf(
			"{\"paymentUid\": \"%s\", \"status\": \"PAID\", \"hotel_id\": %d, \"startDate\": \"%s\", \"endDate\": \"%s\"}",
			payment_uid, checkhotel.Id, ownrequest.Start_date, ownrequest.End_date))
		createreservationbodyReader := bytes.NewReader(createreservationjsonBody)

		createreservationurl := reservationService + "/api/v1/reservations"
		createreservationreq, err := http.NewRequest(http.MethodPost, createreservationurl, createreservationbodyReader)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		createreservationreq.Header.Set("X-User-Name", username)

		createreservationres, err := http.DefaultClient.Do(createreservationreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if createreservationres.StatusCode != http.StatusCreated {
			c.AbortWithStatus(createreservationres.StatusCode)
			return
		}

		temp = strings.Split(createreservationres.Header.Get("Location"), "/")
		reservation_uid := temp[len(temp)-1]

		c.JSON(http.StatusOK, &CreateReservationResponse{
			Reservation_uid: reservation_uid,
			Hotel_uid:       ownrequest.Hotel_uid,
			Start_date:      ownrequest.Start_date,
			End_date:        ownrequest.End_date,
			Payment: &PaymentResponse{
				Price:  price,
				Status: "PAID",
			},
			Status:   "PAID",
			Discount: incrementcounterdata.Discount,
		})

	}
}

func (g *GatewayHandlers) DeleteReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		reservation_uid := c.Param("reservationUid")
		if reservation_uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		getreservationurl := fmt.Sprintf("%s/api/v1/reservations/%s", reservationService, reservation_uid)
		getreservationreq, err := http.NewRequest(http.MethodGet, getreservationurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getreservationreq.Header.Set("X-User-Name", username)

		getreservationres, err := http.DefaultClient.Do(getreservationreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if getreservationres.StatusCode != http.StatusOK {
			c.AbortWithStatus(getreservationres.StatusCode)
			return
		}

		getreservationresBody, err := io.ReadAll(getreservationres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var reservation GetReservationByUid
		if json.Unmarshal(getreservationresBody, &reservation) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		deletereservationurl := fmt.Sprintf("%s/api/v1/reservations/%s", reservationService, reservation_uid)
		deletereservationreq, err := http.NewRequest(http.MethodDelete, deletereservationurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		deletereservationres, err := http.DefaultClient.Do(deletereservationreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if deletereservationres.StatusCode != http.StatusNoContent {
			c.AbortWithStatus(deletereservationres.StatusCode)
			return
		}

		cancelpaymentjsonBody := []byte("{\"status\": \"CANCELED\"}")
		cancelpaymentjsonBodyReader := bytes.NewReader(cancelpaymentjsonBody)

		cancelpaymenturl := fmt.Sprintf("%s/api/v1/payments/%s", paymentService, reservation.Payment_uid)
		cancelpaymentreq, err := http.NewRequest(http.MethodPatch, cancelpaymenturl, cancelpaymentjsonBodyReader)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cancelpaymentres, err := http.DefaultClient.Do(cancelpaymentreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if cancelpaymentres.StatusCode != http.StatusOK {
			c.AbortWithStatus(cancelpaymentres.StatusCode)
			return
		}

		// decrement reservation count
		decrementcounterjsonBody := []byte("{\"reservationCount\": -1}")
		decrementcounterbodyReader := bytes.NewReader(decrementcounterjsonBody)

		decrementcounterurl := loylatyService + "/api/v1/loyalty"
		decrementcounterreq, err := http.NewRequest(http.MethodPatch, decrementcounterurl, decrementcounterbodyReader)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		decrementcounterreq.Header.Set("X-User-Name", username)

		decrementcounterres, err := http.DefaultClient.Do(decrementcounterreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if decrementcounterres.StatusCode != http.StatusOK {
			c.AbortWithStatus(decrementcounterres.StatusCode)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (g *GatewayHandlers) GetReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		reservation_uid := c.Param("reservationUid")
		if reservation_uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		getreservationurl := fmt.Sprintf("%s/api/v1/reservations/%s", reservationService, reservation_uid)
		getreservationreq, err := http.NewRequest(http.MethodGet, getreservationurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getreservationreq.Header.Set("X-User-Name", username)

		getreservationres, err := http.DefaultClient.Do(getreservationreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if getreservationres.StatusCode != http.StatusOK {
			c.AbortWithStatus(getreservationres.StatusCode)
			return
		}

		getreservationresBody, err := io.ReadAll(getreservationres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var reservation GetReservationByUid
		if json.Unmarshal(getreservationresBody, &reservation) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getpaymenturl := fmt.Sprintf("%s/api/v1/payments/%s", paymentService, reservation.Payment_uid)
		getpaymentreq, err := http.NewRequest(http.MethodGet, getpaymenturl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getpaymentres, err := http.DefaultClient.Do(getpaymentreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if getpaymentres.StatusCode != http.StatusOK {
			c.AbortWithStatus(getpaymentres.StatusCode)
			return
		}

		getpaymentresBody, err := io.ReadAll(getpaymentres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var payment PaymentResponse
		if json.Unmarshal(getpaymentresBody, &payment) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, GetOneReservation{
			Reservation_uid: reservation_uid,
			Hotel:           *reservation.Hotel,
			Status:          reservation.Status,
			Start_date:      (reservation.Start_date)[:10],
			End_data:        (reservation.End_data)[:10],
			Payment:         &payment,
		})
	}
}

func (g *GatewayHandlers) GetReservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		getreservationsurl := fmt.Sprintf("%s/api/v1/reservations", reservationService)
		getreservationsreq, err := http.NewRequest(http.MethodGet, getreservationsurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getreservationsreq.Header.Set("X-User-Name", username)

		getreservationsres, err := http.DefaultClient.Do(getreservationsreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if getreservationsres.StatusCode != http.StatusOK {
			c.AbortWithStatus(getreservationsres.StatusCode)
			return
		}

		getreservationsresBody, err := io.ReadAll(getreservationsres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var reservations []*GetReservationByUid
		if json.Unmarshal(getreservationsresBody, &reservations) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		response := make([]*GetOneReservation, len(reservations))

		for i, reservation := range reservations {
			getpaymenturl := fmt.Sprintf("%s/api/v1/payments/%s", paymentService, reservation.Payment_uid)
			getpaymentreq, err := http.NewRequest(http.MethodGet, getpaymenturl, nil)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			getpaymentres, err := http.DefaultClient.Do(getpaymentreq)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if getpaymentres.StatusCode != http.StatusOK {
				c.AbortWithStatus(getpaymentres.StatusCode)
				return
			}

			getpaymentresBody, err := io.ReadAll(getpaymentres.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var payment PaymentResponse
			if json.Unmarshal(getpaymentresBody, &payment) != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			response[i] = &GetOneReservation{
				Reservation_uid: reservation.Reservation_uid,
				Hotel:           *reservation.Hotel,
				Status:          reservation.Status,
				Start_date:      (reservation.Start_date)[:10],
				End_data:        (reservation.End_data)[:10],
				Payment:         &payment,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

func (g *GatewayHandlers) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		getreservationsurl := fmt.Sprintf("%s/api/v1/reservations", reservationService)
		getreservationsreq, err := http.NewRequest(http.MethodGet, getreservationsurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getreservationsreq.Header.Set("X-User-Name", username)

		getreservationsres, err := http.DefaultClient.Do(getreservationsreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if getreservationsres.StatusCode != http.StatusOK {
			c.AbortWithStatus(getreservationsres.StatusCode)
			return
		}

		getreservationsresBody, err := io.ReadAll(getreservationsres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var reservations []*GetReservationByUid
		if json.Unmarshal(getreservationsresBody, &reservations) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		response := make([]*GetOneReservation, len(reservations))

		for i, reservation := range reservations {
			getpaymenturl := fmt.Sprintf("%s/api/v1/payments/%s", paymentService, reservation.Payment_uid)
			getpaymentreq, err := http.NewRequest(http.MethodGet, getpaymenturl, nil)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			getpaymentres, err := http.DefaultClient.Do(getpaymentreq)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if getpaymentres.StatusCode != http.StatusOK {
				c.AbortWithStatus(getpaymentres.StatusCode)
				return
			}

			getpaymentresBody, err := io.ReadAll(getpaymentres.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			var payment PaymentResponse
			if json.Unmarshal(getpaymentresBody, &payment) != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			response[i] = &GetOneReservation{
				Reservation_uid: reservation.Reservation_uid,
				Hotel:           *reservation.Hotel,
				Status:          reservation.Status,
				Start_date:      (reservation.Start_date)[:10],
				End_data:        (reservation.End_data)[:10],
				Payment:         &payment,
			}
		}

		getloyaltyurl := fmt.Sprintf("%s/api/v1/loyalty", loylatyService)
		getloyaltyreq, err := http.NewRequest(http.MethodGet, getloyaltyurl, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		getloyaltyreq.Header.Set("X-User-Name", username)

		getloyaltyres, err := http.DefaultClient.Do(getloyaltyreq)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if getloyaltyres.StatusCode != http.StatusOK {
			c.AbortWithStatus(getloyaltyres.StatusCode)
			return
		}

		getloyaltyresBody, err := io.ReadAll(getloyaltyres.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var loyalty UserLoyalty
		if json.Unmarshal(getloyaltyresBody, &loyalty) != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, UserInfo{
			Reservations: response,
			Loyalty:      &loyalty,
		})
	}
}

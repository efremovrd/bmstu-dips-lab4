package http

import (
	"bmstu-dips-lab2/reservation-service/internal/reservation"

	"github.com/gin-gonic/gin"
)

func MapReservationRoutes(reservationGroup *gin.RouterGroup, h reservation.Handlers) {
	reservationGroup.POST("", h.Create())
	reservationGroup.GET("/:reservationUid", h.GetByReservationUid())
	reservationGroup.GET("", h.GetByUsername())
	reservationGroup.DELETE("/:reservationUid", h.Delete())
}

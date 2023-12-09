package http

import (
	"bmstu-dips-lab2/reservation-service/internal/hotel"

	"github.com/gin-gonic/gin"
)

func MapHotelRoutes(hotelGroup *gin.RouterGroup, h hotel.Handlers) {
	hotelGroup.POST("", h.Create())
	hotelGroup.GET("", h.GetAllPaged())
	hotelGroup.GET("/:uid", h.GetByUid())
}

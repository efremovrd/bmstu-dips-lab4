package http

import (
	"bmstu-dips-lab2/loyalty-service/internal/loyalty"

	"github.com/gin-gonic/gin"
)

func MapPaymentRoutes(loyaltyGroup *gin.RouterGroup, h loyalty.Handlers) {
	loyaltyGroup.POST("", h.Create())
	loyaltyGroup.PATCH("", h.UpdateResCountByOne())
	loyaltyGroup.GET("", h.GetByUsername())
}

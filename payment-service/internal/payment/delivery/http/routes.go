package http

import (
	"bmstu-dips-lab2/payment-service/internal/payment"

	"github.com/gin-gonic/gin"
)

func MapPaymentRoutes(paymentGroup *gin.RouterGroup, h payment.Handlers) {
	paymentGroup.POST("", h.Create())
	paymentGroup.PATCH("/:paymentUid", h.Update())
	paymentGroup.GET("/:paymentUid", h.GetByPaymentUid())
}

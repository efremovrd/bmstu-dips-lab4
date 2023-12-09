package http

import (
	"bmstu-dips-lab2/gateway-service/internal/gateway"

	"github.com/gin-gonic/gin"
)

func MapGatewayRoutes(gatewayGroup *gin.RouterGroup, g gateway.Handlers) {
	gatewayGroup.POST("/reservations", g.CreateReservation())
	gatewayGroup.GET("/hotels", g.GetHotels())
	gatewayGroup.GET("/me", g.GetUserInfo())
	gatewayGroup.GET("/reservations", g.GetReservations())
	gatewayGroup.GET("/reservations/:reservationUid", g.GetReservation())
	gatewayGroup.GET("/loyalty", g.GetLoyalty())
	gatewayGroup.DELETE("/reservations/:reservationUid", g.DeleteReservation())
}

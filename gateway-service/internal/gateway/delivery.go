package gateway

import "github.com/gin-gonic/gin"

type Handlers interface {
	CreateReservation() gin.HandlerFunc
	GetHotels() gin.HandlerFunc
	GetUserInfo() gin.HandlerFunc
	GetReservations() gin.HandlerFunc
	GetReservation() gin.HandlerFunc
	DeleteReservation() gin.HandlerFunc
	GetLoyalty() gin.HandlerFunc
}

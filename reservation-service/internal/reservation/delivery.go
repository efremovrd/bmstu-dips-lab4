package reservation

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create() gin.HandlerFunc
	GetByReservationUid() gin.HandlerFunc
	GetByUsername() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

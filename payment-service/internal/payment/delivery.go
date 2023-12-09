package payment

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	GetByPaymentUid() gin.HandlerFunc
}

package hotel

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create() gin.HandlerFunc
	GetAllPaged() gin.HandlerFunc
	GetByUid() gin.HandlerFunc
}

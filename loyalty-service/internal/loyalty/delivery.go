package loyalty

import "github.com/gin-gonic/gin"

type Handlers interface {
	Create() gin.HandlerFunc
	UpdateResCountByOne() gin.HandlerFunc
	GetByUsername() gin.HandlerFunc
}

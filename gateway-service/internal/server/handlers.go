package server

import (
	h "bmstu-dips-lab2/gateway-service/internal/gateway/delivery/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	gH := h.NewGatewayHandlers()

	s.router.GET("/manage/health", GetHealth())

	api := s.router.Group("/api")

	v1 := api.Group("/v1")

	gateway := v1.Group("")
	h.MapGatewayRoutes(gateway, gH)

	return nil
}

func GetHealth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

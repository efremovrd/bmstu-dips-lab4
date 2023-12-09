package server

import (
	h "bmstu-dips-lab2/payment-service/internal/payment/delivery/http"
	"bmstu-dips-lab2/payment-service/internal/payment/repo"
	"bmstu-dips-lab2/payment-service/internal/payment/usecase"
	googleuuid "bmstu-dips-lab2/pkg/uuider/impl"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	uuider := googleuuid.NewGoogleUUID()

	pRepo := repo.NewPaymentRepo(s.db)
	pUC := usecase.NewPaymentUseCase(pRepo, uuider)
	pH := h.NewPaymentHandlers(pUC)

	s.router.GET("/manage/health", GetHealth())

	api := s.router.Group("/api")

	v1 := api.Group("/v1")

	payments := v1.Group("/payments")
	h.MapPaymentRoutes(payments, pH)

	return nil
}

func GetHealth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

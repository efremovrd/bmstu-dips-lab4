package server

import (
	h "bmstu-dips-lab2/loyalty-service/internal/loyalty/delivery/http"
	"bmstu-dips-lab2/loyalty-service/internal/loyalty/repo"
	"bmstu-dips-lab2/loyalty-service/internal/loyalty/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	lRepo := repo.NewLoyaltyRepo(s.db)
	lUC := usecase.NewLoyaltyUseCase(lRepo)
	lH := h.NewLoyaltyHandlers(lUC)

	s.router.GET("/manage/health", GetHealth())

	api := s.router.Group("/api")

	v1 := api.Group("/v1")

	loyalties := v1.Group("/loyalty")
	h.MapPaymentRoutes(loyalties, lH)

	return nil
}

func GetHealth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

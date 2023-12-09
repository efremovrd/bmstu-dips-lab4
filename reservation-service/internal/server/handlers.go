package server

import (
	h "bmstu-dips-lab2/reservation-service/internal/hotel/delivery/http"
	"bmstu-dips-lab2/reservation-service/internal/hotel/repo"
	"bmstu-dips-lab2/reservation-service/internal/hotel/usecase"
	"net/http"

	h2 "bmstu-dips-lab2/reservation-service/internal/reservation/delivery/http"
	repoR "bmstu-dips-lab2/reservation-service/internal/reservation/repo"
	usecaseR "bmstu-dips-lab2/reservation-service/internal/reservation/usecase"

	googleuuid "bmstu-dips-lab2/pkg/uuider/impl"

	"github.com/gin-gonic/gin"
)

func (s *Server) MapHandlers() error {
	uuider := googleuuid.NewGoogleUUID()

	hRepo := repo.NewHotelRepo(s.db)
	hUC := usecase.NewHotelUseCase(hRepo, uuider)
	hH := h.NewHotelHandlers(hUC)

	rRepo := repoR.NewReservationRepo(s.db)
	rUC := usecaseR.NewReservationUseCase(rRepo, uuider)
	rH := h2.NewReservationHandlers(rUC, hUC)

	s.router.GET("/manage/health", GetHealth())

	api := s.router.Group("/api")

	v1 := api.Group("/v1")

	hotels := v1.Group("/hotels")
	h.MapHotelRoutes(hotels, hH)

	reservations := v1.Group("/reservations")
	h2.MapReservationRoutes(reservations, rH)

	return nil
}

func GetHealth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

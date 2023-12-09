package http

import (
	"bmstu-dips-lab2/loyalty-service/internal/loyalty"
	"bmstu-dips-lab2/loyalty-service/models"
	"bmstu-dips-lab2/pkg/errs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoyaltyCreatRequest struct {
	Username string `json:"username" binding:"required"`
}

type LoyaltyUpdResCountByOneRequest struct {
	Reservation_count *int `json:"reservationCount"`
}

type LoyaltyResponse struct {
	Id                int    `json:"id"`
	Username          string `json:"username"`
	Status            string `json:"status"`
	Reservation_count int    `json:"reservationCount"`
	Discount          int    `json:"discount"`
}

type LoyaltyHandlers struct {
	loyaltyUC loyalty.UseCase
}

func NewLoyaltyHandlers(loyaltyUC loyalty.UseCase) loyalty.Handlers {
	return &LoyaltyHandlers{
		loyaltyUC: loyaltyUC,
	}
}

func (l *LoyaltyHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(LoyaltyCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := LoyaltyCreatRequestToBL(request)

		createdloyalty, err := l.loyaltyUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Header("Location", "/api/v1/loyalty/"+strconv.Itoa(createdloyalty.Id))

		c.Status(http.StatusCreated)
	}
}

func (l *LoyaltyHandlers) UpdateResCountByOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(LoyaltyUpdResCountByOneRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := LoyaltyUpdResCountByOneRequestToBL(request)
		modelBL.Username = username

		updatedloyalty, err := l.loyaltyUC.UpdateResCountByOne(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, LoyaltyBLToResponse(updatedloyalty))
	}
}

func (l *LoyaltyHandlers) GetByUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("X-User-Name")
		if username == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundloyalty, err := l.loyaltyUC.GetByUsername(c, username)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, LoyaltyBLToResponse(foundloyalty))
	}
}

func LoyaltyCreatRequestToBL(dto *LoyaltyCreatRequest) *models.Loyalty {
	return &models.Loyalty{
		Username: dto.Username,
	}
}

func LoyaltyUpdResCountByOneRequestToBL(dto *LoyaltyUpdResCountByOneRequest) *models.Loyalty {
	return &models.Loyalty{
		Reservation_count: *dto.Reservation_count,
	}
}

func LoyaltyBLToResponse(modelBL *models.Loyalty) *LoyaltyResponse {
	return &LoyaltyResponse{
		Id:                modelBL.Id,
		Username:          modelBL.Username,
		Status:            modelBL.Status,
		Reservation_count: modelBL.Reservation_count,
		Discount:          modelBL.Discount,
	}
}

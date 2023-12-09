package http

import (
	"bmstu-dips-lab2/pkg/errs"
	"bmstu-dips-lab2/pkg/types"
	"bmstu-dips-lab2/reservation-service/internal/hotel"
	"bmstu-dips-lab2/reservation-service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelCreatRequest struct {
	Price     int    `json:"price" binding:"required"`
	Stars     int    `json:"stars" binding:"required"`
	Name      string `json:"name" binding:"required"`
	City      string `json:"city" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Hotel_uid string `json:"hotelUid"`
}

type HotelResponse struct {
	Id        int    `json:"id"`
	Price     int    `json:"price"`
	Stars     int    `json:"stars"`
	Name      string `json:"name"`
	Hotel_uid string `json:"hotelUid"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Country   string `json:"country"`
}

type HotelsResponse struct {
	Items         []*HotelResponse `json:"items"`
	TotalElements int              `json:"totalElements"`
	Page          int              `json:"page"`
	PageSize      int              `json:"pageSize"`
}

type HotelHandlers struct {
	hotelUC hotel.UseCase
}

func NewHotelHandlers(hotelUC hotel.UseCase) hotel.Handlers {
	return &HotelHandlers{
		hotelUC: hotelUC,
	}
}

func (p *HotelHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(HotelCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := HotelCreatRequestToBL(request)

		createdhotel, err := p.hotelUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Header("Location", "/api/v1/hotels/"+strconv.Itoa(createdhotel.Id))

		c.Status(http.StatusCreated)
	}
}

func (h *HotelHandlers) GetAllPaged() gin.HandlerFunc {
	return func(c *gin.Context) {
		intpage, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		intsize, err := strconv.Atoi(c.Query("size"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		limit, offset, ok := types.ValidateGetSets(c.Query("size"), strconv.Itoa((intpage-1)*intsize))
		if !ok {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundhotels, err := h.hotelUC.GetAllPaged(c, types.GetSets{Limit: limit, Offset: offset})
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, HotelsBLToResponse(foundhotels, intpage, intsize))
	}
}

func (h *HotelHandlers) GetByUid() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")
		if uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundhotel, err := h.hotelUC.GetByUid(c, uid)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, HotelBLToResponse(foundhotel))
	}
}

func HotelCreatRequestToBL(dto *HotelCreatRequest) *models.Hotel {
	return &models.Hotel{
		Stars:     dto.Stars,
		Price:     dto.Price,
		Name:      dto.Name,
		City:      dto.City,
		Address:   dto.Address,
		Country:   dto.Country,
		Hotel_uid: dto.Hotel_uid,
	}
}

func HotelBLToResponse(model *models.Hotel) *HotelResponse {
	return &HotelResponse{
		Id:        model.Id,
		Price:     model.Price,
		Stars:     model.Stars,
		Name:      model.Name,
		Address:   model.Address,
		Country:   model.Country,
		City:      model.City,
		Hotel_uid: model.Hotel_uid,
	}
}

func HotelsBLToResponse(hotels []*models.Hotel, page, size int) *HotelsResponse {
	if hotels == nil {
		return nil
	}

	items := make([]*HotelResponse, len(hotels))

	for i, p := range hotels {
		items[i] = &HotelResponse{
			Stars:     p.Stars,
			Price:     p.Price,
			Id:        p.Id,
			City:      p.City,
			Country:   p.Country,
			Address:   p.Address,
			Name:      p.Name,
			Hotel_uid: p.Hotel_uid,
		}
	}

	return &HotelsResponse{
		Page:          page,
		PageSize:      size,
		TotalElements: len(items),
		Items:         items,
	}
}

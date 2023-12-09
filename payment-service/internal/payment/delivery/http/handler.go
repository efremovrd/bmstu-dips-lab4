package http

import (
	"bmstu-dips-lab2/payment-service/internal/payment"
	"bmstu-dips-lab2/payment-service/models"
	"bmstu-dips-lab2/pkg/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentCreatRequest struct {
	Price  int    `json:"price" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type PaymentUpdRequest struct {
	Status *string `json:"status"`
	Price  *int    `json:"price"`
}

type PaymentResponse struct {
	Id          int    `json:"id"`
	Price       int    `json:"price"`
	Status      string `json:"status"`
	Payment_uid string `json:"paymentUid"`
}

type PaymentHandlers struct {
	paymentUC payment.UseCase
}

func NewPaymentHandlers(paymentUC payment.UseCase) payment.Handlers {
	return &PaymentHandlers{
		paymentUC: paymentUC,
	}
}

func (p *PaymentHandlers) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(PaymentCreatRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		modelBL := PaymentCreatRequestToBL(request)

		createdpayment, err := p.paymentUC.Create(c, modelBL)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.Header("Location", "/api/v1/payments/"+createdpayment.Payment_uid)

		c.Status(http.StatusCreated)
	}
}

func (p *PaymentHandlers) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := new(PaymentUpdRequest)

		err := c.ShouldBindJSON(request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		uid := c.Param("paymentUid")
		if uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		const updStr string = "y"
		const updInt int = 1

		modelBL := &models.Payment{Payment_uid: uid}
		toUpdate := &models.Payment{}

		if request.Status != nil {
			modelBL.Status = *request.Status
			toUpdate.Status = updStr
		}

		if request.Price != nil {
			modelBL.Price = *request.Price
			toUpdate.Price = updInt
		}

		updatedpayment, err := p.paymentUC.Update(c, modelBL, toUpdate)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, PaymentBLToResponse(updatedpayment))
	}
}

func (p *PaymentHandlers) GetByPaymentUid() gin.HandlerFunc {
	return func(c *gin.Context) {
		payment_uid := c.Param("paymentUid")
		if payment_uid == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		foundpayment, err := p.paymentUC.GetByPaymentUid(c, payment_uid)
		if err != nil {
			c.AbortWithStatus(errs.MatchHttpErr(err))
			return
		}

		c.JSON(http.StatusOK, PaymentBLToResponse(foundpayment))
	}
}

func PaymentCreatRequestToBL(dto *PaymentCreatRequest) *models.Payment {
	return &models.Payment{
		Status: dto.Status,
		Price:  dto.Price,
	}
}

func PaymentBLToResponse(modelBL *models.Payment) *PaymentResponse {
	return &PaymentResponse{
		Id:          modelBL.Id,
		Price:       modelBL.Price,
		Status:      modelBL.Status,
		Payment_uid: modelBL.Payment_uid,
	}
}

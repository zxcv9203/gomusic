package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zxcv9203/gomusic/backend/src/dblayer"
)

type HandlerInterface interface {
	GetProducts(c *gin.Context)
	GetPromos(c *gin.Context)
	AddUser(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	GetOrders(c *gin.Context)
	Charge(c *gin.Context)
}

type Handler struct {
	db dblayer.DBLayer
}

func NewHandler() (*Handler, error) {
	// Handler 객체에 대한 포인터 생성
	return new(Handler), nil
}

func (h *Handler) GetProducts(c *gin.Context) {
	if h.db == nil {
		return
	}
	products, err := h.db.GetAllProducts()
	if err != nil {
		/*
			첫 번째 인자는 HTTP 상태 코드, 두 번째는 응답의 바디
		*/
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusOK, products)
}

package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zxcv9203/gomusic/backend/src/dblayer"
	"github.com/zxcv9203/gomusic/backend/src/models"
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

// Handler 객체에 대한 포인터를 생성하는 함수
func NewHandler() (*Handler, error) {
	// Handler 객체에 대한 포인터 생성
	return new(Handler), nil
}

// 상품목록 조회
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

// 프로모션 조회
func (h *Handler) GetPromos(c *gin.Context) {
	if h.db == nil {
		return
	}
	promos, err := h.db.GetPromos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, promos)
}

// 사용자 로그인
func (h *Handler) SignIn(c *gin.Context) {
	if h.db == nil {
		return
	}
	var customer models.Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, err = h.db.SignInUser(customer.Email, customer.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// 사용자 가입
func (h *Handler) AddUser(c *gin.Context) {
	if h.db == nil {
		return
	}
	var customer models.Customer
	err := c.ShouldBindJSON(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, err = h.db.AddUser(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// 사용자 로그아웃
func (h *Handler) SignOut(c *gin.Context) {
	if h.db == nil {
		return
	}
	p := c.Param("id")
	// id는 현재 문자형이기 때문에 정수형으로 변환해야 합니다.
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.db.SignOutUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// 사용자 주문 내역 조회
func (h *Handler) GetOrders(c *gin.Context) {
	if h.db == nil {
		return
	}
	// id 매개변수 추출
	p := c.Param("id")
	// p 를 string -> int 변경
	id, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 데이터베이스 레이어 메서드 호출과 주문 내역 조회
	orders, err := h.db.GetCustomerOrdersById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) Charge(c *gin.Context) {
	if h.db == nil {
		return
	}
}

func MyCustomMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 요청을 처리하기 전에 실행할 코드
		// 예제 변수 설정
		c.Set("v", "123")
		// c.Get("v")를 하면 변수 값을 확인할 수 있습니다.

		// 요청 처리 로직 실행
		c.Next()

		// 이 코드는 핸들러 실행이 끝나면 실행됩니다.

		// 응답 코드 확인
		status := c.Writer.Status()
		//status를 활용하는 코드 추가
	}
}

// MyCustomLogger 요청 처리후 전후에 특정 문자열을 출력하는 간단한 미들웨어
func MyCustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("*******************************************")
		c.Next()
		fmt.Println("*******************************************")
	}
}
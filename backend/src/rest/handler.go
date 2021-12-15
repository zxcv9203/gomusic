package rest

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
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
func NewHandler(dbtype, conn string) (HandlerInterface, error) {
	db, err := dblayer.NewORM(dbtype, conn)
	if err != nil {
		return nil, err
	}
	return &Handler{
		db: db,
	}, nil
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
		if err == dblayer.ErrINVALIDPASSWORD {
			// 잘못된 패스워드의 경우 forbidden http 에러 반환
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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
	err = h.db.SignOutUserByID(id)
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
	orders, err := h.db.GetCustomerOrdersByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) Charge(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server database error"})
		return
	}
	// 프론트엔드에서 전달받은 데이터를 나타내는 구조체
	request := struct {
		models.Order
		Remember    bool   `json:"rememberCard"`
		UseExisting bool   `json:"useExisting"`
		Token       string `json:"token"`
	}{}
	err := c.ShouldBindJSON(&request)
	// 파싱 중 에러 발생시 에러 보고 후 반환
	if err != nil {
		c.JSON(http.StatusBadRequest, request)
		return
	}
	// 스트라이프 API 키 설정(테스트 키)
	stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"
	// *stripe.ChargeParams 타입 인스턴스 생성
	chargeP := &stripe.ChargeParams{
		// 요청에 명시된 판매 가격
		Amount: stripe.Int64(int64(request.Price)),
		// 결제 통화
		Currency: stripe.String("usd"),
		// 설명
		Description: stripe.String("GoMusic charge..."),
	}
	// 스트라이프 사용자 ID 초기화
	stripeCustomerID := ""
	//이미 저장해둔 신용카드로 결제하는 경우인지 확인
	if request.UseExisting {
		//저장된 카드 사용
		log.Println("Getting credit card id....")
		// 스트라이프 사용자 ID를 데이터베이스에서 조회하는 메서드
		stripeCustomerID, err = h.db.GetCreditCardCID(request.CustomerId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		/*
			새로운 신용카드로 결제하는 경우 *stripe.CustomerParams 타입 인스턴스를 생성하고
			이를 사용해 *stripe.Customer 타입 인스턴스를 생성합니다.
		*/
		cp := &stripe.CustomerParams{}
		cp.SetSource(request.Token)
		customer, err := customer.New(cp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stripeCustomerID = customer.ID
		if request.Remember {
			// 스트라이프 사용자 id를 저장하고 데이터베이스에 저장된 사용자 ID와 연결합니다.
			err = h.db.SaveCreditCardForCustomer(request.CustomerId, stripeCustomerID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
	/* 동일 상품 주문 여부 확인 없이 새로운 주문으로 가정 */
	// *stripe.ChargeParams 타입 인스턴스에 스트라이프 사용자 ID를 설정합니다.
	chargeP.Customer = stripe.String(stripeCustomerID)
	// 신용카드 결제 요청
	_, err = charge.New(chargeP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 주문 내용 데이터베이스에 저장
	err = h.db.AddOrder(request.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

//func MyCustomMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 요청을 처리하기 전에 실행할 코드
//		// 예제 변수 설정
//		c.Set("v", "123")
//		// c.Get("v")를 하면 변수 값을 확인할 수 있습니다.

//		// 요청 처리 로직 실행
//		c.Next()

//		// 이 코드는 핸들러 실행이 끝나면 실행됩니다.

//		// 응답 코드 확인
//		status := c.Writer.Status()
//		//status를 활용하는 코드 추가
//	}
//}

// MyCustomLogger 요청 처리전에 특정 문자열을 출력하는 간단한 미들웨어
func MyCustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("*******************************************")
		c.Next()
		fmt.Println("*******************************************")
	}
}

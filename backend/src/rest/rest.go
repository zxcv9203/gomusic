package rest

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RunAPIWithHandler(address string, h HandlerInterface) error {
	log.Println("handler")
	// Gin 엔진(기본 미들웨어 미사용)
	r := gin.New()
	// 요청 처리 전후로 특정 문자열을 출력하는 미들웨어
	r.Use(MyCustomLogger())
	// 상품 목록
	r.GET("/products", h.GetProducts)
	// 프로모션 목록
	r.GET("/promos", h.GetPromos)
	userGroup := r.Group("/user")
	{
		// 사용자 로그아웃
		userGroup.POST("/:id/signout", h.SignOut)
		// 주문 내역
		userGroup.GET("/:id/orders", h.GetOrders)
	}
	usersGroup := r.Group("/users")
	{
		// 사용자 로그인
		usersGroup.POST("/signin", h.SignIn)
		// 사용자 추가
		usersGroup.POST("", h.AddUser)
		// 결제
		usersGroup.POST("charge", h.Charge)
	}
	// 서버 시작
	return (r.Run(address))
}

func RunAPI(address string) error {
	h, err := NewHandler()
	if err != nil {
		return err
	}
	log.Println("runapi")
	return RunAPIWithHandler(address, h)
}

package rest

import (
	"github.com/gin-gonic/gin"
)

func RunAPI(address string) error {
	r := gin.Default()
	//상품 목록 응답
	r.GET("/products", func(c *gin.Context) {
		// 클라이언트에게 상품 목록 반환
	})
	//프로모션 목록 응답
	r.GET("promos", func(c *gin.Context) {
		// 클라이언트에게 프로모션 목록 반환
	})
	// 사용자 로그인 POST 요청
	r.POST("/users/signin", func(c *gin.Context) {
		//사용자 로그인
	})
	// 사용자 가입 POST 요청
	r.POST("/users", func(c *gin.Context) {
		//사용자 가입
	})
	// 사용자 로그아웃 POST 요청
	/*
		아래 경로는 사용자 ID를 포함해야 합니다.
		ID는 사용자마다 고유한 값이기 때문에 와일드 카드(*)를 사용하며 ":id"는 id라는 이름의 변수를 의미합니다.
	*/
	r.POST("/users/:id/signout", func(c *gin.Context) {
		// 해당 ID의 사용자 로그아웃
	})
	//구매 목록 조회
	r.GET("/user/:id/orders", func(c *gin.Context) {
		// 해당 ID의 사용자의 주문 내역 조회
	})
	// 결제 POST 요청
	r.POST("/users/charge", func(c *gin.Context) {
		// 신용카드 결제 처리
	})
	return nil
}

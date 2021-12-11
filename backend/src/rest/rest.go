package rest

import (
	"github.com/gin-gonic/gin"
)

func RunAPI(address string) error {
	r := gin.Default()
	r.GET("/products", func(c *gin.Context) {
		// 클라이언트에게 상품 목록 반환
	})
}

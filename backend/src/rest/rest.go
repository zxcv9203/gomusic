package rest

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// socket 연결을 위한 데이터를 세팅해두는 구조체
type ConnPath struct {
	username   string
	password   string
	socketPath string
	database   string
}

func RunAPIWithHandler(address string, h HandlerInterface) error {
	// Gin 엔진(기본 미들웨어 미사용)
	//r := gin.New()
	// Gin 엔진(기본 미들웨어 사용)
	r := gin.Default()
	// 요청 처리 전후로 특정 문자열을 출력하는 미들웨어
	r.Use(static.ServeRoot("/", "../../public"))
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
	return (r.RunTLS(address, "../../cert.pem", "../../key.pem"))
}

// Mysql을 Scoket으로 연결하기 위해 사용하는 함수
func GetSocketConn() string {
	connInfo := ConnPath{
		username:   "root",
		password:   "Rladydcjf12!",
		socketPath: "/var/run/mysqld/mysqld.sock",
		database:   "gomusic",
	}
	conn := connInfo.username + ":" + connInfo.password +
		"@unix(" + connInfo.socketPath + ")" + "/" +
		connInfo.database + "?charset=utf8"
	return conn
}

func RunAPI(address string) error {
	h, err := NewHandler("mysql", GetSocketConn())
	if err != nil {
		return err
	}
	return RunAPIWithHandler(address, h)
}

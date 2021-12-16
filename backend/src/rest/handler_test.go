package rest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zxcv9203/gomusic/backend/src/dblayer"
)

type errMSG struct {
	Error string `json:"error"`
}

/*
	GetProducts 함수는 gomusic에서 판매하는 모든 상품의 목록을 반환하는 HTTP 핸들러 함수입니다.
	HTTP 핸들러 함수는 특정 URL로 HTTP 요청이 들어오면 호출되는 함수이며, 핸들러는 HTTP 요청을 처리하고 HTTP를 통해서 응답을 보냅니다.
*/
func TestHandler_GetProducts(t *testing.T) {
	// 테스트 모드 활성화로 로깅 방지
	gin.SetMode(gin.TestMode)
	// 데이터를 하드 코딩한 모의 객체로 테스트
	mockdbLayer := dblayer.NewMockDBLayerWithData()
	h := NewHandlerWithDB(mockdbLayer)
	const productsURL string = "/products"
	// 테스트 케이스를 저장하는 구조체 슬라이스
	tests := []struct {
		name             string
		inErr            error
		outStatusCode    int
		expectedRespBody interface{}
	}{
		{
			"getproductsnoerrors",
			nil,
			http.StatusOK,
			mockdbLayer.GetMockProductData(),
		},
		{
			"getproductswitcherror",
			errors.New("get products error"),
			http.StatusInternalServerError,
			errMSG{Error: "get products error"},
		},
	}

}

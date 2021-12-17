package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zxcv9203/gomusic/backend/src/dblayer"
	"github.com/zxcv9203/gomusic/backend/src/models"
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
	for _, tt := range tests {
		// 서브테스트 정의
		t.Run(tt.name, func(t *testing.T) {
			// 서브테스트 실행
			mockdbLayer.SetError(tt.inErr)
			// 테스트 요청 생성
			req := httptest.NewRequest(http.MethodGet, productsURL, nil)
			//  HTTP response recoder 생성
			w := httptest.NewRecorder()
			// response recoder를 사용해 gin 엔진 객체를 생성합니다.
			// context 인스턴스는 사용하지 않습니다.
			_, engine := gin.CreateTestContext(w)
			// Gin 엔진 인스턴스를 사용해 productsURL에 GetProducts를 매핑합니다.
			engine.GET(productsURL, h.GetProducts)
			// Gin 엔진이 HTTP 요청을 처리하도록 설정하고 HTTP 응답은 ResponseRecorder 타입으로 생성합니다.
			engine.ServeHTTP(w, req)
			// 결과 검증
			response := w.Result()
			// 전달받은 HTTP 코드와 세팅한 HTTP 코드가 다를경우 에러처리
			if response.StatusCode != tt.outStatusCode {
				t.Errorf("Received Status code %d does not match expected status code %d", response.StatusCode, tt.outStatusCode)
			}
			// http 응답의 형식을 미리 알 수 없기 때문에 interface{} 타입을 사용합니다.
			var respBody interface{}
			// 에러가 발생한 경우 응답을 errMSG 타입으로 변환
			if tt.inErr != nil {
				var errmsg errMSG
				json.NewDecoder(response.Body).Decode(&errmsg)
				// 에러 메시지를 respBody에 저장
				respBody = errmsg
			} else {
				// 에러가 없을 경우 응답을 product 타입의 슬라이스로 변환
				products := []models.Product{}
				json.NewDecoder(response.Body).Decode(&products)
				// 디코딩한 상품 목록을 respBody에 저장
			}
			if !reflect.DeepEqual(respBody, tt.expectedRespBody) {
				t.Errorf("Received HTTP response body %+v does not match expected HTTP response Body %+v", respBody, tt.expectedRespBody)
			}
		})

	}
}

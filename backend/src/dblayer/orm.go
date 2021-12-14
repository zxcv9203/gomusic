package dblayer

import (
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zxcv9203/gomusic/backend/src/models"
	"golang.org/x/crypto/bcrypt"
)

type DBORM struct {
	*gorm.DB
}

// NewORM DB 종류와 연결주소를 전달받아 연결하는 생성자 함수입니다.
func NewORM(dbname, con string) (*DBORM, error) {
	db, err := gorm.Open(dbname, con)
	return &DBORM{
		DB: db,
	}, err
}

// GetAllProducts 모든 상품의 목록을 반환하는 메서드입니다.
func (db *DBORM) GetAllProducts() (products []models.Product, err error) {
	return products, db.Find(&products).Error
}

// GetPromos 현재 프로모션 중인 상품을 반환하는 메서드입니다.
func (db *DBORM) GetPromos() (products []models.Product, err error) {
	return products, db.Where("promotion IS NOT NULL").Find(&products).Error
}

// GetCustomerByName 사용자 이름과 성을 인자로 전달받고 사용자 정보를 반환하는 메서드입니다.
func (db *DBORM) GetCustomerByName(firstname string, lastname string) (customer models.Customer, err error) {
	return customer, db.Where(&models.Customer{FirstName: firstname, LastName: lastname}).Find(&customer).Error
}

// GetCustomerByID 사용자 ID로 사용자 정보를 반환하는 메서드입니다.
func (db *DBORM) GetCustomerByID(id int) (customer models.Customer, err error) {
	return customer, db.First(&customer, id).Error
}

// GetProduct ID가 가리키는 상품의 정보를 반환하는 메서드입니다.
func (db *DBORM) GetProduct(id int) (product models.Product, err error) {
	return product, db.First(&product, id).Error
}

// AddUser 새로운 사용자의 정보를 데이터베이스에 삽입하는 메서드입니다.
func (db *DBORM) AddUser(customer models.Customer) (models.Customer, error) {
	/*
		이 메서드는 사용자의 패스워드를 해싱하고 로그인 상태로 설정합니다.
	*/
	hashPassword(&customer.Pass)
	customer.LoggedIn = true
	// gorm Create() 메서드는 테이블에 데이터를 삽입합니다.
	err := db.Create(&customer).Error
	customer.Pass = ""
	return customer, err
}

// SignInUser 사이트 로그인을 하는 함수
func (db *DBORM) SignInUser(email, pass string) (customer models.Customer, err error) {
	// Where에 해당하는 사용자 행을 가져옴
	result := db.Table("customers").Where(&models.Customer{Email: email})
	// 이메일로 사용자 정보 조회
	err = result.First(&customer).Error
	if err != nil {
		return customer, err
	}
	// 패스워드 문자열과 해시 값 비교
	if !checkPassword(customer.Pass, pass) {
		return customer, ErrINVALIDPASSWORD
	}
	// 패스워드 문자열이 공유되지 않도록 비교를 완료하면 지웁니다.
	customer.Pass = ""
	// loggedin 필드 업데이트
	err = result.Update("loggedin", 1).Error
	if err != nil {
		return customer, err
	}
	// 사용자 행 반환
	return customer, result.Find(&customer).Error
}

// SignOutUserByID 사이트 로그아웃 하는 메서드
func (db *DBORM) SignOutUserByID(id int) error {
	//ID에 해당하는 사용자 구조체 생성
	customer := models.Customer{
		Model: gorm.Model{
			ID: uint(id),
		},
	}
	// 사용자의 상태를 로그아웃 상태로 업데이트하는 함수
	return db.Table("customers").Where(&customer).Update("loggedin", 0).Error
}

// GetCustomerOrdersByID 특정 사용자의 주문 내역을 조회하는 메서드
func (db *DBORM) GetCustomerOrdersByID(id int) (orders []models.Order, err error) {
	/*
		orders와 customers, products 테이블을 조인한 후 customers 테이블에서 전달받은 id 값에 해당하는 사용자 정보를 조회합니다.
		그리고 products 테이블에서 현재 선택된 상품 ID에 해당하는 상품 정보를 가져옵니다.
	*/
	return orders, db.Table("orders").Select("*").
		Joins("join customers on customers.id = customer_id").
		Joins("join products on products.id = product_id").
		Where("customer_id=?", id).Scan(&orders).Error
}

func hashPassword(s *string) error {
	if s == nil {
		return errors.New("Reference provided for hashing password is nil")
	}
	// bcrypt 패키지에서 사용할 수 있게 패스워드 문자열을 바이트 슬라이스로 변환합니다.
	sBytes := []byte(*s)
	// GenerateFromPassword() 메서드는 패스워드 해시를 반환합니다.
	hashBytes, err := bcrypt.GenerateFromPassword(sBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 패스워드 문자열을 해시 값으로 바꿉니다.
	*s = string(hashBytes[:])
	return nil
}

// 해시와 패스워드 문자열이 일치하는지 확인하는 메서드입니다.
func checkPassword(existingHash, incomingPass string) bool {
	/*
		해시와 패스워드 문자열이 일치하지 않으면
		아래 메서드는 에러를 반환합니다.
	*/
	return bcrypt.CompareHashAndPassword([]byte(existingHash),
		[]byte(incomingPass)) == nil
}

// orders 테이블에 결제 내역 추가
func (db *DBORM) AddOrder(order models.Order) error {
	return db.Create(&order).Error
}

// 신용카드 ID 조회
func (db *DBORM) GetCreditCardCID(id int) (string, error) {
	customerWithCCID := struct {
		models.Customer
		CCID string `gorm:"column:cc_customerid"`
	}{}
	return customerWithCCID.CCID, db.First(&customerWithCCID, id).Error
}

// 신용카드 정보 저장
func (db *DBORM) SaveCreditCardForCustomer(id int, ccid string) error {
	result := db.Table("customers").Where("id=?", id)
	return result.Update("cc_customerid", ccid).Error
}

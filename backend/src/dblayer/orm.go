package dblayer

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zxcv9203/gomusic/backend/src/models"
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
	return customer, db.Create(&customer).Error
}

// SignInUser 사이트 로그인을 하는 함수
func (db *DBORM) SignInUser(email, pass string) (customer models.Customer, err error) {
	if !checkPassword(pass) {
		return customer, errors.New("Invalid password")
	}
	// Where에 해당하는 사용자 행을 가져옴
	result := db.Table("customers").Where(&models.Customer{Email: email})
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
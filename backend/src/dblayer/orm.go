package dblayer

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/zxcv9203/gomusic/backend/src/models"
)

type DBORM struct {
	*gorm.DB
}

// DB 종류와 연결주소를 전달받아 연결하는 생성자 함수입니다.
func NewORM(dbname, con string) (*DBORM, error) {
	db, err := gorm.Open(dbname, con)
	return &DBORM{
		DB: db,
	}, err
}

// 모든 상품의 목록을 반환하는 메서드입니다.
func (db *DBORM) GetAllProducts() (products []models.Product, err error) {
	return products, db.Find(&products).Error
}

// 현재 프로모션 중인 상품을 반환하는 메서드입니다.
func (db *DBORM) GetPromos() (products []models.Product, err error) {
	return products, db.Where("promotion IS NOT NULL").Find(&products).Error
}

// 사용자 이름과 성을 인자로 전달받고 사용자 정보를 반환하는 메서드입니다.
func (db *DBORM) GetCustomerByName(firstname string, lastname string) (customer models.Customer, err error) {
	return customer, db.Where(&models.Customer{FirstName: firstname, LastName: lastname}).Find(&customer).Error
}

// 사용자 ID로 사용자 정보를 반환하는 메서드입니다.
func (db *DBORM) GetCustomerByID(id int) (customer models.Customer, err error) {
	return customer, db.First(&customer, id).Error
}

// ID가 가리키는 상품의 정보를 반환하는 메서드입니다.
func (db *DBORM) GetProduct(id int) (product models.Product, err error) {
	return product, db.First(&product, id).Error
}

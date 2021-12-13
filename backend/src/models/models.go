package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Image       string  `json:"img"`
	ImagAlt     string  `gorm:"column:imgalt" json:"imgalt"`
	Price       float64 `json:"price"`
	Promotion   float64 `json:"promotion"` // sql.NullFloat64
	ProductName string  `gorm:"column:productname" json:"productname"`
	Description string
}

// TableName을 따로 지정하지않으면 첫번째 글자가 소문자가 되고
// 마지막 글자에 s를 추가시키고 해당 이름의 테이블과 매핑합니다.
// 구조체 : Product -> 테이블 : products 매핑
func (Product) TableName() string {
	return "products"
}

type Customer struct {
	gorm.Model
	FirstName string  `gorm:"column:firstname" json:"firstname"`
	LastName  string  `gorm:"column:lastname" json:"lastname"`
	Email     string  `gorm:"column:email" json:"email"`
	Pass      string  `json:"password"`
	LoggedIn  bool    `gorm:"column:loggedin" json:"loggedin"`
	Orders    []Order `json:"orders"`
}

type Order struct {
	gorm.Model
	Product
	Customer
	CustomerId   int       `gorm:"column:customer_id"`
	ProductId    int       `gorm:"column:product_id"`
	Price        float64   `gorm:"column:price" json:"sell_price"`
	PurchaseDate time.Time `gorm:"column:purchase_date" json:"purchase_date"`
}

func (Order) TableName() string {
	return "orders"
}

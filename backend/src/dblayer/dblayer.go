package dblayer

import "github.com/zxcv9203/gomusic/backend/src/models"

type DBLayer interface {
	GetAllProducts() ([]models.Product, error)
}

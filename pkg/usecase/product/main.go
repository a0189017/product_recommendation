package product

import (
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/repository/product"

	"github.com/redis/go-redis/v9"
)

type ProductUsecase struct {
	db                database.DBInterface
	redis             *redis.Client
	productRepository product.ProductRepository
}

func NewProductUsecase(db database.DBInterface, redis *redis.Client) *ProductUsecase {
	return &ProductUsecase{
		productRepository: product.NewProductRepository(db.GetDB(), redis),
	}
}

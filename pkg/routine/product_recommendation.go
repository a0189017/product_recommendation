package routine

import (
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/usecase/product"

	"github.com/redis/go-redis/v9"
)

func ProductRecommendation(database database.DBInterface, redis *redis.Client) error {
	uc := product.NewProductUsecase(database, redis)
	err := uc.UpdateProductRecommendation()

	return err
}

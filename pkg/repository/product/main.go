package product

import (
	"context"
	"encoding/json"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProductRecommendation() (productRecommendationList []*model.Product, err error)
	GetProductRecommendationFromRedis() ([]*model.Product, error)
	SetProductRecommendationToRedis(productRecommendationList []*model.Product) error
}

type productRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProductRepository(db *gorm.DB, redis *redis.Client) ProductRepository {
	return &productRepository{db: db, redis: redis}
}

func (r *productRepository) GetProductRecommendation() (productRecommendationList []*model.Product, err error) {
	// get product recommendation from db
	// for testing, we only get 10 products
	err = r.db.Find(&productRecommendationList).Limit(10).Error
	return
}

func (r *productRepository) GetProductRecommendationFromRedis() ([]*model.Product, error) {
	productRecommendation, err := r.redis.Get(context.Background(), constants.ProductRecommendationKey).Result()
	if err != nil {
		return nil, err
	}

	productRecommendationList := []*model.Product{}
	err = json.Unmarshal([]byte(productRecommendation), &productRecommendationList)
	if err != nil {
		return nil, err
	}

	return productRecommendationList, nil
}

func (r *productRepository) SetProductRecommendationToRedis(productRecommendationList []*model.Product) error {
	productRecommendation, err := json.Marshal(productRecommendationList)
	if err != nil {
		return err
	}

	return r.redis.Set(context.Background(), constants.ProductRecommendationKey, productRecommendation, 10*time.Minute).Err()
}

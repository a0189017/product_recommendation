package product

import (
	"errors"
	"product_recommendation/mocks"
	"product_recommendation/pkg/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_UpdateProductRecommendation(t *testing.T) {
	mockRepo := new(mocks.ProductRepository)
	uc := &ProductUsecase{productRepository: mockRepo}

	productId := uuid.New()
	mockProduct := &model.Product{
		BaseModel:   &model.BaseModel{Id: productId},
		Title:       "Test Product",
		Price:       100,
		Description: "Test Description",
		Category:    "Test Category",
	}
	productList := []*model.Product{mockProduct}

	t.Run("[POSITIVE] update recommendation success", func(t *testing.T) {
		mockRepo.On("GetProductRecommendation").Return(productList, nil).Once()
		mockRepo.On("SetProductRecommendationToRedis", productList).Return(nil).Once()
		err := uc.UpdateProductRecommendation()
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("[NEGATIVE] get product recommendation failed", func(t *testing.T) {
		mockRepo.On("GetProductRecommendation").Return(nil, errors.New("db error")).Once()
		err := uc.UpdateProductRecommendation()
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("[NEGATIVE] set product recommendation to redis failed", func(t *testing.T) {
		mockRepo.On("GetProductRecommendation").Return(productList, nil).Once()
		mockRepo.On("SetProductRecommendationToRedis", productList).Return(errors.New("redis error")).Once()
		err := uc.UpdateProductRecommendation()
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

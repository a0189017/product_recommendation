package product

import (
	"errors"
	"testing"

	"product_recommendation/mocks"
	"product_recommendation/pkg/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetProductRecommendation(t *testing.T) {
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

	t.Run("[POSITIVE] get recommendation success from redis", func(t *testing.T) {
		mockRepo.On("GetProductRecommendationFromRedis").Return(productList, nil).Once()
		result, err := uc.GetProductRecommendation()
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, productId.String(), result[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("[POSITIVE] redis miss, db success", func(t *testing.T) {
		mockRepo.On("GetProductRecommendationFromRedis").Return(nil, errors.New("redis miss")).Once()
		mockRepo.On("GetProductRecommendation").Return(productList, nil).Once()
		mockRepo.On("SetProductRecommendationToRedis", productList).Return(nil).Once()
		result, err := uc.GetProductRecommendation()
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, productId.String(), result[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("[NEGATIVE] redis miss, db error", func(t *testing.T) {
		mockRepo.On("GetProductRecommendationFromRedis").Return(nil, errors.New("redis miss")).Once()
		mockRepo.On("GetProductRecommendation").Return(nil, errors.New("db error")).Once()
		result, err := uc.GetProductRecommendation()
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("[NEGATIVE] redis miss, db success but set redis error", func(t *testing.T) {
		mockRepo.On("GetProductRecommendationFromRedis").Return(nil, errors.New("redis miss")).Once()
		mockRepo.On("GetProductRecommendation").Return(productList, nil).Once()
		mockRepo.On("SetProductRecommendationToRedis", productList).Return(errors.New("set redis error")).Once()
		result, err := uc.GetProductRecommendation()
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, productId.String(), result[0].ID)
		mockRepo.AssertExpectations(t)
	})
}

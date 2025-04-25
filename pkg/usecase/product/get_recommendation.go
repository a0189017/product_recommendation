package product

import (
	"product_recommendation/pkg/model"
	"product_recommendation/pkg/schema"
)

func (uc *ProductUsecase) GetProductRecommendation() ([]*schema.ProductListResponse, error) {
	var productRecommendationList []*model.Product
	productRecommendationList, err := uc.productRepository.GetProductRecommendationFromRedis()

	// if product recommendation list is not found in redis, get it from db
	if err != nil {
		productRecommendationList, err = uc.productRepository.GetProductRecommendation()
		if err != nil {
			return nil, err
		}

		err = uc.productRepository.SetProductRecommendationToRedis(productRecommendationList)
	}

	productListResponse := []*schema.ProductListResponse{}
	for _, product := range productRecommendationList {
		productListResponse = append(productListResponse, &schema.ProductListResponse{
			ID:          product.Id.String(),
			Title:       product.Title,
			Price:       product.Price,
			Description: product.Description,
			Category:    product.Category,
		})
	}

	return productListResponse, nil
}

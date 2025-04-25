package product

func (uc *ProductUsecase) UpdateProductRecommendation() error {
	productRecommendationList, err := uc.productRepository.GetProductRecommendation()
	if err != nil {
		return err
	}

	return uc.productRepository.SetProductRecommendationToRedis(productRecommendationList)
}

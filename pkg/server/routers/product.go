package routers

import (
	"net/http"
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/usecase/product"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetProductRecommendation(c *gin.Context) {
	handler := func(c *gin.Context, database database.DBInterface, redis *redis.Client) error {
		// TODO: get product recommendations with optional filters?
		// TODO: like different user, different recommendation
		productUsecase := product.NewProductUsecase(database, redis)
		productList, err := productUsecase.GetProductRecommendation()
		if err != nil {
			return err
		}

		response := GenResponse(productList)
		c.JSON(http.StatusOK, response)
		return nil
	}

	WrapperWithRedis(c, handler)
}

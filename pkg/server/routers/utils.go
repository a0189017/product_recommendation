package routers

import (
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GenResponse(data interface{}) interface{} {
	return gin.H{"data": data}
}

func Wrapper(c *gin.Context, handlerFunc func(c *gin.Context, d database.DBInterface) error) {
	d := c.MustGet(constants.FieldDatabase)

	err := handlerFunc(c, d.(database.DBInterface))
	if err != nil {
		systemError := ParseSystemError(c, err)
		c.JSON(systemError.StatusCode, gin.H{"error": systemError})
		c.Abort()
	} else {
		c.Next()
	}
}

func WrapperWithRedis(c *gin.Context, handlerFunc func(c *gin.Context, d database.DBInterface, r *redis.Client) error) {
	d := c.MustGet(constants.FieldDatabase)
	r := c.MustGet(constants.FieldRedis)

	err := handlerFunc(c, d.(database.DBInterface), r.(*redis.Client))
	if err != nil {
		systemError := ParseSystemError(c, err)
		c.JSON(systemError.StatusCode, gin.H{"error": systemError})
		c.Abort()
	} else {
		c.Next()
	}
}

func WrapperErr(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			systemError := ParseSystemError(c, err)
			c.JSON(systemError.StatusCode, gin.H{"error": systemError})
			c.Abort()
		}
	}
}

func ParseSystemError(c *gin.Context, err error) (systemError errors.SystemError) {
	if h, ok := err.(errors.SystemError); ok {
		systemError = h
	} else {
		systemError = errors.New(errors.ErrorInfo{Err: err})
	}
	c.Set(constants.StackTrace, systemError.StackTrace)
	return
}

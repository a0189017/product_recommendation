package routers

import (
	"net/http"
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/errors"
	"product_recommendation/pkg/schema"
	"product_recommendation/pkg/types"
	"product_recommendation/pkg/usecase/auth"
	"product_recommendation/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	handler := func(c *gin.Context, database database.DBInterface) error {
		var request schema.AuthRegisterRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			return errors.ErrorUserPasswordFormatInvalid()
		}
		uc := auth.NewAuthUsecase(database)
		if err := uc.Register(request); err != nil {
			return err
		}

		// also response user data
		response := GenResponse(types.H{
			"success": true,
		})

		c.JSON(http.StatusOK, response)
		return nil
	}

	Wrapper(c, handler)
}

func Login(c *gin.Context) {
	handler := func(c *gin.Context, database database.DBInterface) error {
		var request schema.AuthLoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			return errors.ErrorUserPasswordFormatInvalid()
		}

		uc := auth.NewAuthUsecase(database)
		token, otp, err := uc.Login(request)
		if err != nil {
			return err
		}

		response := GenResponse(schema.AuthLoginResponse{
			Success: true,
			Otp:     otp,
			Token:   token,
		})

		c.JSON(http.StatusOK, response)
		return nil
	}

	Wrapper(c, handler)
}

func VerifyOTP(c *gin.Context) {
	handler := func(c *gin.Context, database database.DBInterface) error {
		loginUser := utils.GetLoginUser(c)
		var request schema.AuthVerifyOTPRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			return errors.ErrorInvalidOTP()
		}
		request.Account = loginUser.Account

		uc := auth.NewAuthUsecase(database)
		token, err := uc.VerifyOTP(request)
		if err != nil {
			return err
		}

		response := GenResponse(schema.AuthVerifyOTPResponse{
			Success: true,
			Token:   token,
		})

		c.JSON(http.StatusOK, response)
		return nil
	}

	Wrapper(c, handler)
}

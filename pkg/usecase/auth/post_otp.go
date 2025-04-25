package auth

import (
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/errors"
	"product_recommendation/pkg/schema"
	"product_recommendation/pkg/types"
	"product_recommendation/pkg/utils"
	"time"
)

func (uc *AuthUsecase) VerifyOTP(request schema.AuthVerifyOTPRequest) (token string, err error) {
	if request.Otp == nil {
		err = errors.ErrorMissingOTP()
		return
	}
	user, err := uc.authRepository.GetUser(request.Account)
	if err != nil {
		err = errors.ErrorUserNotFound()
		return
	}

	if user.OtpExpired != nil && user.OtpExpired.Time().Before(time.Now()) {
		err = errors.ErrorOTPExpired()
		return
	}

	isMatch := utils.ComparePassword(*user.Otp, *request.Otp)
	if !isMatch {
		err = errors.ErrorInvalidOTP()
		return
	}

	if user.IsOtpAuth {
		err = errors.ErrorOTPAlreadyUsed()
		return
	}

	// generate token
	now := time.Now().UTC()
	accessTokenExpiresIn := config.GetConfig().Settings.AccessTokenExpireSec
	accessToken, err := utils.GenerateToken(utils.TokenPayload{
		GateUserClaims: types.GateUserClaims{
			UserId:    user.Id,
			Account:   user.Account,
			LoginType: constants.LoginTypeToken,
		},
		Now:       now,
		ExpiresIn: accessTokenExpiresIn,
	})
	if err != nil {
		return
	}

	// update user is otp auth
	user.IsOtpAuth = true
	err = uc.authRepository.UpdateUser(user)
	if err != nil {
		return
	}

	return accessToken, nil
}

package auth

import (
	"crypto/rand"
	"math/big"
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/errors"
	"product_recommendation/pkg/schema"
	"product_recommendation/pkg/types"
	"product_recommendation/pkg/utils"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (uc *AuthUsecase) Login(request schema.AuthLoginRequest) (token string, otp string, err error) {
	config := config.GetConfig()

	user, err := uc.authRepository.GetUser(request.Account)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if err == gorm.ErrRecordNotFound {
		err = errors.ErrorUserNotFound()
		return
	}

	// check password
	if !utils.ComparePassword(user.Credential, request.Password) {
		err = errors.ErrorInvalidPassword()
		return
	}

	// generate otp
	otpLength := config.Settings.OtpLength
	otp, err = generateOTP(otpLength)
	if err != nil {
		return
	}
	// hash otp
	otpHash, err := utils.GeneratePasswordHash(otp)
	if err != nil {
		return
	}

	user.Otp = &otpHash
	expireTime := time.Now().Add(5 * time.Minute).UTC()
	user.OtpExpired = lo.ToPtr(types.Iso8601Time(expireTime))
	user.IsOtpAuth = false

	err = uc.authRepository.UpdateUser(user)
	if err != nil {
		return
	}

	now := time.Now().UTC()
	accessTokenExpiresIn := "300"
	accessToken, err := utils.GenerateToken(utils.TokenPayload{
		GateUserClaims: types.GateUserClaims{
			UserId:    user.Id,
			Account:   user.Account,
			LoginType: constants.LoginTypeOtp,
		},
		Now:       now,
		ExpiresIn: accessTokenExpiresIn,
	})
	if err != nil {
		return
	}

	return accessToken, otp, utils.SendEmail(user.Account, utils.EmailTemplate{
		Subject: "OTP Code",
		Body:    otp,
	})
}

func generateOTP(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[nBig.Int64()]
	}
	return string(result), nil
}

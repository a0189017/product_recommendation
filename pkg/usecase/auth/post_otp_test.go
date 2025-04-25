package auth

import (
	"errors"
	"testing"
	"time"

	"product_recommendation/mocks"
	"product_recommendation/pkg/model"
	"product_recommendation/pkg/schema"
	"product_recommendation/pkg/types"
	"product_recommendation/pkg/utils"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_PostOtp(t *testing.T) {
	mockRepo := new(mocks.AuthRepository)
	uc := &AuthUsecase{authRepository: mockRepo}
	baseModel := &model.BaseModel{Id: uuid.New()}
	otp := "123456"
	hashOtp, _ := utils.GeneratePasswordHash(otp)
	now := types.Iso8601Time(time.Now().UTC().Add(10 * time.Minute))

	t.Run("[POSITIVE] verify otp success", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  baseModel,
			Account:    "test@test.com",
			Otp:        &hashOtp,
			OtpExpired: &now,
			IsOtpAuth:  false,
		}
		mockRepo.On("GetUser", "test@test.com").Return(user, nil)
		mockRepo.On("UpdateUser", mock.Anything).Return(nil)
		patches := gomonkey.ApplyFunc(utils.ComparePassword, func(hashedPassword string, password string) bool {
			return true
		})
		defer patches.Reset()
		patchesToken := gomonkey.ApplyFunc(utils.GenerateToken, func(payload utils.TokenPayload) (string, error) {
			return "token", nil
		})
		defer patchesToken.Reset()

		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "test@test.com",
			Otp:     &otp,
		})
		assert.NoError(t, err)
		assert.Equal(t, "token", token)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] missing otp", func(t *testing.T) {
		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "test@test.com",
			Otp:     nil,
		})
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("[NEGATIVE] user not found", func(t *testing.T) {
		mockRepo.On("GetUser", "notfound@test.com").Return(nil, errors.New("not found"))
		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "notfound@test.com",
			Otp:     &otp,
		})
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] otp expired", func(t *testing.T) {
		expired := types.Iso8601Time(time.Now().UTC().Add(-10 * time.Second))
		user := &model.Login{
			BaseModel:  baseModel,
			Account:    "test@test.com",
			Otp:        &hashOtp,
			OtpExpired: &expired,
			IsOtpAuth:  false,
		}
		mockRepo.On("GetUser", "test@test.com").Return(user, nil)
		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "test@test.com",
			Otp:     &otp,
		})
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] invalid otp", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  baseModel,
			Account:    "test@test.com",
			Otp:        &hashOtp,
			OtpExpired: &now,
			IsOtpAuth:  false,
		}
		mockRepo.On("GetUser", "test@test.com").Return(user, nil)
		patches := gomonkey.ApplyFunc(utils.ComparePassword, func(hashedPassword string, password string) bool {
			return false
		})
		defer patches.Reset()
		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "test@test.com",
			Otp:     &otp,
		})
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] otp already used", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  baseModel,
			Account:    "test@test.com",
			Otp:        &hashOtp,
			OtpExpired: &now,
			IsOtpAuth:  true,
		}
		mockRepo.On("GetUser", "test@test.com").Return(user, nil)
		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "test@test.com",
			Otp:     &otp,
		})
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] token generate fail", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  baseModel,
			Account:    "test@test.com",
			Otp:        &hashOtp,
			OtpExpired: &now,
			IsOtpAuth:  false,
		}
		mockRepo.On("GetUser", "test@test.com").Return(user, nil)
		patches := gomonkey.ApplyFunc(utils.ComparePassword, func(hashedPassword string, password string) bool {
			return true
		})
		defer patches.Reset()
		patchesToken := gomonkey.ApplyFunc(utils.GenerateToken, func(payload utils.TokenPayload) (string, error) {
			return "", errors.New("token error")
		})
		defer patchesToken.Reset()
		token, err := uc.VerifyOTP(schema.AuthVerifyOTPRequest{
			Account: "test@test.com",
			Otp:     &otp,
		})
		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})
}

package auth

import (
	"errors"
	"testing"

	"product_recommendation/mocks"
	"product_recommendation/pkg/model"
	"product_recommendation/pkg/schema"
	"product_recommendation/pkg/utils"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gorm.io/gorm"
)

func Test_PostLogin(t *testing.T) {
	mockRepo := new(mocks.AuthRepository)
	uc := &AuthUsecase{authRepository: mockRepo}
	hashPw, _ := utils.GeneratePasswordHash("A1234567b+")
	BaseModel := &model.BaseModel{
		Id: uuid.New(),
	}

	t.Run("[POSITIVE] login success", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  BaseModel,
			Account:    "test@test.com",
			Credential: hashPw,
			IsActive:   true,
		}
		mockRepo.On("GetUser", "test@test.com").Return(user, nil)
		mockRepo.On("UpdateUser", mock.Anything).Return(nil)

		token, otp, err := uc.Login(schema.AuthLoginRequest{
			Account:  "test@test.com",
			Password: "A1234567b+",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, otp)
		mockRepo.AssertExpectations(t)

		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] user not found", func(t *testing.T) {
		mockRepo.On("GetUser", "notfound@test.com").Return(nil, gorm.ErrRecordNotFound)
		_, _, err := uc.Login(schema.AuthLoginRequest{
			Account:  "notfound@test.com",
			Password: "A1234567b+",
		})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)

		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] db error", func(t *testing.T) {
		mockRepo.On("GetUser", "dberr@test.com").Return(nil, errors.New("db error"))
		_, _, err := uc.Login(schema.AuthLoginRequest{
			Account:  "dberr@test.com",
			Password: "A1234567b+",
		})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)

		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] invalid password", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  BaseModel,
			Account:    "wrongpw@test.com",
			Credential: hashPw,
			IsActive:   true,
		}
		mockRepo.On("GetUser", "wrongpw@test.com").Return(user, nil)
		// 需 monkey patch utils.ComparePassword 回傳 false
		_, _, err := uc.Login(schema.AuthLoginRequest{
			Account:  "wrongpw@test.com",
			Password: "wrongpw",
		})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)

		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] update user fail", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  BaseModel,
			Account:    "failupdate@test.com",
			Credential: hashPw,
			IsActive:   true,
		}
		mockRepo.On("GetUser", "failupdate@test.com").Return(user, nil)
		mockRepo.On("UpdateUser", mock.Anything).Return(errors.New("update error"))
		_, _, err := uc.Login(schema.AuthLoginRequest{
			Account:  "failupdate@test.com",
			Password: "A1234567b+",
		})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)

		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] token fail", func(t *testing.T) {
		user := &model.Login{
			BaseModel:  BaseModel,
			Account:    "failtoken",
			Credential: hashPw,
			IsActive:   true,
		}
		mockRepo.On("GetUser", "failtoken").Return(user, nil)
		mockRepo.On("UpdateUser", mock.Anything).Return(nil)

		patches := gomonkey.ApplyFunc(utils.GenerateToken, func(payload utils.TokenPayload) (string, error) {
			return "", errors.New("token error")
		})
		defer patches.Reset()

		_, _, err := uc.Login(schema.AuthLoginRequest{
			Account:  "failtoken",
			Password: "A1234567b+",
		})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)

		ClearMock(mockRepo, nil)
	})
}

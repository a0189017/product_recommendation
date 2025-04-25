package auth

import (
	"testing"

	"product_recommendation/mocks"
	"product_recommendation/pkg/model"
	"product_recommendation/pkg/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func Test_PostRegister(t *testing.T) {
	mockRepo := new(mocks.AuthRepository)
	uc := &AuthUsecase{authRepository: mockRepo}

	t.Run("[POSITIVE] register success", func(t *testing.T) {
		mockRepo.On("GetUser", "test@test.com").Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("CreateUser", mock.Anything).Return(nil)
		req := schema.AuthRegisterRequest{
			Account:  "test@test.com",
			Password: "A1234567b+",
		}
		err := uc.Register(req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] user already exists", func(t *testing.T) {
		mockRepo.On("GetUser", "exist@test.com").Return(&model.Login{Account: "exist@test.com"}, nil)
		req := schema.AuthRegisterRequest{
			Account:  "exist@test.com",
			Password: "A1234567b+",
		}
		err := uc.Register(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "User already exists")
		mockRepo.AssertExpectations(t)
		ClearMock(mockRepo, nil)
	})

	t.Run("[NEGATIVE] password format invalid", func(t *testing.T) {
		req := schema.AuthRegisterRequest{
			Account:  "test2@test.com",
			Password: "123",
		}
		err := uc.Register(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "User password format invalid")
	})

	t.Run("[NEGATIVE] account format invalid", func(t *testing.T) {
		req := schema.AuthRegisterRequest{
			Account:  "notanemail",
			Password: "A1234567b+",
		}
		err := uc.Register(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "User account format invalid")
	})
}

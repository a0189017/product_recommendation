package auth

import (
	"product_recommendation/pkg/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *model.Login) error
	GetUser(account string) (*model.Login, error)
	UpdateUser(user *model.Login) error
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

type authRepository struct {
	db *gorm.DB
}

func (r *authRepository) CreateUser(user *model.Login) error {
	return r.db.Create(user).Error
}

func (r *authRepository) GetUser(account string) (*model.Login, error) {
	var user model.Login
	if err := r.db.Where("account = ?", account).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(user *model.Login) error {
	return r.db.Save(user).Error
}

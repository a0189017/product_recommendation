package auth

import (
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/repository/auth"
)

type AuthUsecase struct {
	db             database.DBInterface
	authRepository auth.AuthRepository
}

func NewAuthUsecase(db database.DBInterface) *AuthUsecase {
	return &AuthUsecase{authRepository: auth.NewAuthRepository(db.GetDB())}
}

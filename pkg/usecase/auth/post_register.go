package auth

import (
	"product_recommendation/pkg/errors"
	"product_recommendation/pkg/model"
	"product_recommendation/pkg/schema"
	"product_recommendation/pkg/utils"
	"regexp"
	"unicode"

	"gorm.io/gorm"
)

func (uc *AuthUsecase) Register(request schema.AuthRegisterRequest) error {
	// verify account
	if err := VerifyAccount(request.Account); err != nil {
		return err
	}
	// verify password
	if err := VerifyPassword(request.Password); err != nil {
		return err
	}
	// check if user already exists
	findUser, err := uc.authRepository.GetUser(request.Account)
	if findUser != nil {
		return errors.ErrorUserAlreadyExists()
	}
	// create user
	user := &model.Login{
		Account:  request.Account,
		IsActive: true,
	}

	if err == gorm.ErrRecordNotFound {
		hashedPassword, err := utils.GeneratePasswordHash(request.Password)
		user.Credential = hashedPassword
		if err != nil {
			return err
		}

		err = uc.authRepository.CreateUser(user)
		if err != nil {
			return err
		}
	}

	return nil
}

func VerifyPassword(s string) error {
	var sixOrMore, sixteenOrLess, lower, upper, special bool
	for _, c := range s {
		switch {
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	sixOrMore = len(s) >= 6
	sixteenOrLess = len(s) <= 16

	if sixOrMore && sixteenOrLess && lower && upper && special {
		return nil
	}

	return errors.ErrorUserPasswordFormatInvalid()
}

func VerifyAccount(s string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(s) {
		return errors.ErrorUserAccountFormatInvalid()
	}
	return nil
}

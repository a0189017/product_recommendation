package utils

import (
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/types"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type EmailTemplate struct {
	Subject string
	Body    string
}

func SendEmail(email string, emailTemplate EmailTemplate) error {
	// In this case, we need to send an email containing the OTP
	// This function should support sending emails with different templates based on the use case
	return nil
}

func GeneratePasswordHash(password string) (hashedPassword string, err error) {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hashedPassword = string(hash)
	return
}

func ComparePassword(hashedPassword string, password string) bool {
	byteHash := []byte(hashedPassword)
	bytePwd := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}

type TokenPayload struct {
	types.GateUserClaims
	ExpiresIn string
	Now       time.Time
}

func GenerateToken(payload TokenPayload) (string, error) {
	expireAt, err := ExpireTime(payload.Now, payload.ExpiresIn)
	if err != nil {
		return "", err
	}
	// set claims and sign
	claims := types.GateUserClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  payload.Account,
			ExpiresAt: expireAt.Unix(),
			IssuedAt:  payload.Now.Unix(),
		},
		UserId:    payload.UserId,
		LoginType: payload.LoginType,
		Account:   payload.Account,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(config.GetConfig().Settings.JWTSecret))
	return token, err
}

func ExpireTime(now time.Time, expireSec string) (time.Time, error) {
	expiresIn, err := strconv.Atoi(expireSec)
	if err != nil {
		return time.Time{}, err
	}
	expiresInDuration := time.Duration(expiresIn) * time.Second
	return now.Add(expiresInDuration), nil
}

func GetLoginUser(c *gin.Context) *types.LoginUser {
	u, _ := c.Get(constants.FieldLoginUser)
	if u != nil {
		return u.(*types.LoginUser)
	}
	return nil
}

package types

import (
	"database/sql/driver"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type H map[string]interface{}

type Iso8601Time time.Time

func (i Iso8601Time) Value() (driver.Value, error) {
	return i.Time(), nil
}

func (i Iso8601Time) Time() time.Time {
	var t time.Time = time.Time(i)
	return t
}

type LoginType string

type GateUserClaims struct {
	jwt.StandardClaims
	UserId    uuid.UUID `json:"user_id"`
	LoginType LoginType `json:"login_type"`
	Account   string    `json:"account"`
}

type LoginUser struct {
	UserId    uuid.UUID `json:"user_id"`
	LoginType LoginType `json:"login_type"`
	Account   string    `json:"account"`
}

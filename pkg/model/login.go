package model

import (
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/types"

	"gorm.io/gorm"
)

type Login struct {
	*BaseModel
	Account    string             `json:"account" gorm:"type:varchar(255);uniqueIndex:login_unique_index;not null"`
	Credential string             `json:"credential,omitempty" gorm:"type:text;" swaggerignore:"true"`
	IsActive   bool               `json:"is_active" gorm:"type:boolean;not null;"`
	LastLogin  *types.Iso8601Time `json:"last_login" gorm:"type:timestamp;"`
	Otp        *string            `json:"otp" gorm:"type:varchar(255);"`
	OtpExpired *types.Iso8601Time `json:"otp_expired" gorm:"type:timestamp;"`
	IsOtpAuth  bool               `json:"is_otp_auth" gorm:"type:boolean;not null;"`
}

func (m *Login) TableName() string {
	return "login"
}

func (m *Login) BeforeCreate(tx *gorm.DB) (err error) {
	if m.BaseModel == nil {
		m.BaseModel = &BaseModel{}
	}
	m.attachId()
	m.attachLoginUserId(tx)
	return
}

func (m *Login) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.BaseModel == nil {
		m.BaseModel = &BaseModel{}
	}
	m.attachLoginUserId(tx)
	return
}

func (m *Login) BeforeDelete(tx *gorm.DB) (err error) {
	if m.BaseModel == nil {
		m.BaseModel = &BaseModel{}
	}
	m.attachLoginUserId(tx)
	return
}

func (m *Login) AfterCreate(tx *gorm.DB) (err error) {
	createLog(tx, m, constants.RecordCreate)
	return
}

func (m *Login) AfterUpdate(tx *gorm.DB) (err error) {
	createLog(tx, m, constants.RecordUpdate)
	return
}

func (m *Login) AfterDelete(tx *gorm.DB) (err error) {
	createLog(tx, m, constants.RecordDelete)
	return
}

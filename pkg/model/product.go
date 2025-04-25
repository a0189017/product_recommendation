package model

import (
	"product_recommendation/pkg/constants"

	"gorm.io/gorm"
)

type Product struct {
	*BaseModel
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func (m *Product) TableName() string {
	return "product"
}

func (m *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if m.BaseModel == nil {
		m.BaseModel = &BaseModel{}
	}
	m.attachId()
	m.attachLoginUserId(tx)
	return
}

func (m *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.BaseModel == nil {
		m.BaseModel = &BaseModel{}
	}
	m.attachLoginUserId(tx)
	return
}

func (m *Product) AfterCreate(tx *gorm.DB) (err error) {
	createLog(tx, m, constants.RecordCreate)
	return
}

func (m *Product) AfterUpdate(tx *gorm.DB) (err error) {
	createLog(tx, m, constants.RecordUpdate)
	return
}

func (m *Product) AfterDelete(tx *gorm.DB) (err error) {
	createLog(tx, m, constants.RecordDelete)
	return
}

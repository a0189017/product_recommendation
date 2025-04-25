package migration

import (
	"product_recommendation/pkg/types"
	"product_recommendation/pkg/utils/gormigrate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id         uuid.UUID         `json:"id" gorm:"type:varchar(36);primarykey;"`
	LastUserId string            `json:"-"  gorm:"type:varchar(255);not null;"`
	CreatedAt  types.Iso8601Time `json:"created_at"  gorm:"type:timestamp;not null;"`
	UpdatedAt  types.Iso8601Time `json:"updated_at"  gorm:"type:timestamp;not null;"`
}

func (m *BaseModel) GetId() uuid.UUID {
	if m == nil {
		return uuid.Nil
	} else {
		return m.Id
	}
}

func (m *BaseModel) GetLastUserId() string {
	if m == nil {
		return ""
	} else {
		return m.LastUserId
	}
}

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

var Migration = gormigrate.Migration{
	ID: "00_add_demo_product",
	Migrate: func(db *gorm.DB) error {
		db.AutoMigrate(&Product{})
		products := []*Product{
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "無線藍牙耳機",
				Price:       1299,
				Description: "高音質無線藍牙耳機，續航時間長達12小時",
				Category:    "電子產品",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "智能手錶",
				Price:       2999,
				Description: "多功能智能手錶，支持心率監測和GPS定位",
				Category:    "電子產品",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "咖啡機",
				Price:       4999,
				Description: "全自動咖啡機，可製作多種咖啡飲品",
				Category:    "家電",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "空氣清淨機",
				Price:       5999,
				Description: "高效能空氣清淨機，適用於大坪數空間",
				Category:    "家電",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "登山背包",
				Price:       1999,
				Description: "防水耐磨登山背包，容量50L",
				Category:    "戶外用品",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "電動滑板車",
				Price:       8999,
				Description: "折疊式電動滑板車，續航里程30公里",
				Category:    "交通工具",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "藍牙喇叭",
				Price:       1599,
				Description: "防水藍牙喇叭，360度環繞音效",
				Category:    "電子產品",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "電競鍵盤",
				Price:       2499,
				Description: "機械式電競鍵盤，RGB背光",
				Category:    "電腦周邊",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "電動牙刷",
				Price:       899,
				Description: "聲波電動牙刷，三種清潔模式",
				Category:    "個人護理",
			},
			{
				BaseModel:   &BaseModel{Id: uuid.New()},
				Title:       "行動電源",
				Price:       599,
				Description: "20000mAh大容量行動電源，雙USB輸出",
				Category:    "電子產品",
			},
		}

		for _, product := range products {
			if err := db.Create(product).Error; err != nil {
				return err
			}
		}

		return db.Save(&products).Error
	},
	Rollback: func(db *gorm.DB) error {
		/*
			TODO:
			rollback
		*/
		return nil
	},
}

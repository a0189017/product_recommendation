package model

import (
	"encoding/json"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILogger interface {
	GetId() uuid.UUID
	GetLastUserId() string
	TableName() string
}

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

func (m *BaseModel) attachId() {
	if m.Id == uuid.Nil {
		m.Id = uuid.New()
	}
}

func (m *BaseModel) attachLoginUserId(tx *gorm.DB) {
	if m != nil {
		userId := uuid.Nil.String()
		v := tx.Statement.Context.Value(constants.LoginUserContextKey)
		if v != nil {
			s, ok := v.(string)
			if ok {
				userId = s
			}
		}
		m.LastUserId = userId
	}
}

type TableLog struct {
	*BaseModel
	TableName string    `json:"table_name"  gorm:"type:varchar(36);not null;"`
	RefDataId uuid.UUID `json:"ref_data_id"  gorm:"type:varchar(36);index:idx_table_logs_ref_data_id;not null;"`
	Action    string    `json:"action"  gorm:"type:varchar(6);not null;"`
	Data      string    `json:"data" gorm:"type:json"`
}

func createLog(tx *gorm.DB, logData ILogger, action string) {
	lastUserId := logData.GetLastUserId()
	refDataId := logData.GetId()
	tableName := logData.TableName()

	data, _ := json.Marshal(logData)
	tableLog := &TableLog{
		BaseModel: &BaseModel{
			Id:         uuid.New(),
			LastUserId: lastUserId,
		},
		TableName: tableName,
		Action:    action,
		RefDataId: refDataId,
		Data:      string(data),
	}
	tx.Create(tableLog)
}

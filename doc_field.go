package main

import (
	"time"
)

type LearningEcDocField struct {
	Id        int64     `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`                         // 自增ID
	IsDeleted int       `gorm:"column:is_deleted;default:0;NOT NULL" json:"is_deleted"`                 // 是否删除：0-正常，1-删除
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // 更新时间
}

func (m *LearningEcDocField) TableName() string {
	return "ec_doc_field"
}

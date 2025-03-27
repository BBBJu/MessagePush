package models

import (
	"messagePush/database"
	"time"
)

type Template struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"type:varchar(256);uniqueIndex" json:"name"` // 添加uniqueIndex
	Content    string    `gorm:"type:varchar(4096)" json:"content"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	ModifyTime time.Time `gorm:"autoUpdateTime" json:"modify_time"`
}

func (Template) TableName() string {
	return "template"
}

func (t *Template) GetTemplateByName() error {
	return database.DB.Model(&Template{}).Where("name = ?", t.Name).First(t).Error
}
func (t *Template) GetTemplateById() error {
	return database.DB.Model(&Template{}).Where("id = ?", t.ID).First(t).Error
}

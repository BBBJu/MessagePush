package models

import (
	"messagePush/database"
	"time"
)

type Message struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MsgID        string    `gorm:"type:varchar(256);uniqueIndex" json:"msg_id"`
	TemplateID   int64     `gorm:"type:bigint" json:"template_id"`  // 修改为MySQL支持的类型
	TemplateData string    `gorm:"type:varchar(256)" json:"template_data"`
	SourceID     string    `gorm:"type:varchar(256)" json:"source_id"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"create_time"`
	ModifyTime   time.Time `gorm:"autoUpdateTime" json:"modify_time"`
}

func (Message) TableName() string {
	return "message"
}

func (m *Message) CreateMessage() error {
	return database.DB.Create(m).Error
}

func GetMessageByMsgId(MsgID string) (Message, error) {
	var message Message
	err := database.DB.Where("msg_id = ?", MsgID).First(&message).Error
	return message, err
}

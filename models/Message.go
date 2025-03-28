package models

import (
	"messagePush/database"
	"time"
)

type Message struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MsgID        string    `gorm:"type:varchar(256);uniqueIndex" json:"msg_id"`
	TemplateID   int64     `gorm:"type:bigint" json:"template_id"` // 修改为MySQL支持的类型
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
func BatchGetMessageByMsgIds(MsgIDs []string) ([]Message, error) {

	var messages []Message
	if err := database.DB.Where("msg_id in ?", MsgIDs).Find(&messages).Error; err != nil {
		return nil, err
	}
	// 构建消息指针映射表（节省内存）
	msgMap := make(map[string]*Message, len(messages))
	for i := range messages {
		msg := &messages[i]
		msgMap[msg.MsgID] = msg
	}
	result := make([]Message, len(MsgIDs))
	for i, msgID := range MsgIDs {
		if msg, exists := msgMap[msgID]; exists {
			result[i] = *msg // 复制消息到结果数组
		}
	}
	return result, nil
}

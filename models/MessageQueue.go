package models

import (
	"messagePush/database"
	"time"
)

type MessageQueue struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MsgID        string    `gorm:"type:varchar(256);uniqueIndex" json:"msg_id"`
	To           string    `gorm:"type:varchar(256);index" json:"to"`
	Subject      string    `gorm:"type:varchar(256)" json:"subject"`
	Channel      int       `gorm:"type:int(10)" json:"channel"`
	TemplateID   string    `gorm:"type:varchar(256)" json:"template_id"`
	TemplateData string    `gorm:"type:varchar(256)" json:"template_data"`
	Status       int       `gorm:"type:tinyint(3);default:1" json:"status"`
	Priority     int       `gorm:"type:int(20);index" json:"priority"`
	CreateTime   time.Time `gorm:"autoCreateTime" json:"create_time"`
	ModifyTime   time.Time `gorm:"autoUpdateTime" json:"modify_time"`
}

func (MessageQueue) TableName() string {
	return "message_queue"
}

func (mq *MessageQueue) CreateMessageQueue() error {
	return database.DB.Create(mq).Error
}

func GetMessageQueueByMsgIDs(msgID []string) ([]MessageQueue, error) {
	messageQueues := make([]MessageQueue, len(msgID))
	err := database.DB.Where("msg_id in ?", msgID).Find(&messageQueues).Error
	return messageQueues, err
}

func (mq *MessageQueue) UpdateMessageQueue() error {
	return database.DB.Save(mq).Error
}

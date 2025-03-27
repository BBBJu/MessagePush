package models

import (
	"messagePush/database"
	"time"
)

const (
	MessageStatusCreated = 1
	MessageStatusSending = 2
	MessageStatusSuccess = 3
	MessageStatusFail    = 4
	MessageStatusDead    = 5 //超过重试次数，代表死亡
)

type MessageQueue struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MsgID      string    `gorm:"type:varchar(256);uniqueIndex" json:"msg_id"`
	To         string    `gorm:"type:varchar(256);index" json:"to"`
	Subject    string    `gorm:"type:varchar(256)" json:"subject"`
	Channel    int       `gorm:"type:int(10)" json:"channel"`
	Status     int       `gorm:"type:tinyint(3);default:1" json:"status"`
	Priority   int64     `gorm:"type:int(20);index" json:"priority"`
	OrderBy    int64     `gorm:"type:int(20);index" json:"order"`
	RetryCount int       `gorm:"type:int(10);default:0" json:"retry_count"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	ModifyTime time.Time `gorm:"autoUpdateTime" json:"modify_time"`
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

func GetFailedMessageQueue() ([]MessageQueue, error) {
	messageQueues := make([]MessageQueue, 0)
	err := database.DB.Where("status = ?", MessageStatusFail).Find(&messageQueues).Error
	return messageQueues, err
}

func GetPendingMessages(messageSize int) ([]MessageQueue, error) {
	messageQueues := make([]MessageQueue, messageSize)
	err := database.DB.Where("status in ?", []int{MessageStatusCreated, MessageStatusFail}).
		Order("order_by asc").
		Limit(messageSize).
		Find(&messageQueues).Error
	return messageQueues, err
}

func BatchCreateMessageQueue(messageQueues []MessageQueue) error {
	if len(messageQueues) == 0 {
		return nil
	}
	return database.DB.Create(&messageQueues).Error
}

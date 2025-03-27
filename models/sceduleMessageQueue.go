package models

import (
	"messagePush/database"
	"time"
)

type ScheduleMessageQueue struct {
	ID               int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	MsgID            string `gorm:"type:varchar(256);uniqueIndex" json:"msg_id"`
	To               string `gorm:"type:varchar(256);index" json:"to"`
	Subject          string `gorm:"type:varchar(256)" json:"subject"`
	Channel          int    `gorm:"type:int(10)" json:"channel"`
	ProcessTimeStamp int64  `gorm:"type:bigint" json:"process_time_stamp"`
	//status 只有处理和未处理两种状态，0未处理，1处理。
	Status     int       `gorm:"type:int(10)" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	ModifyTime time.Time `gorm:"autoUpdateTime" json:"modify_time"`
}

func (s *ScheduleMessageQueue) CreateScheduleMessageQueue() error {
	return database.DB.Model(&ScheduleMessageQueue{}).Create(s).Error
}

func GetScheduleMessageQueuesByIds(ids []string) ([]ScheduleMessageQueue, error) {
	var scheduleMessageQueues []ScheduleMessageQueue
	err := database.DB.Model(scheduleMessageQueues).Where("msg_id IN ?", ids).Find(&scheduleMessageQueues).Error
	return scheduleMessageQueues, err
}

func BatchUpdateScheduleMessageQueueStatus(scheduleMessageQueues []ScheduleMessageQueue) error {
	if len(scheduleMessageQueues) == 0 {
		return nil
	}
	return database.DB.Model(&ScheduleMessageQueue{}).Updates(scheduleMessageQueues).Error
}
